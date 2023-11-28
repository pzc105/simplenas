#include <unordered_map>
#include <string_view>
#include "bt_service.hpp"
#include "libtorrent/read_resume_data.hpp"
#include "libtorrent/write_resume_data.hpp"
#include "libtorrent/magnet_uri.hpp"
#include "translate.hpp"
#include "bt_plugin.hpp"

using namespace std;
using namespace grpc;

namespace prpc
{
  bt_service::bt_service() : _pusher_manager(this)
  {
  }

  bt_service::~bt_service()
  {
    if (_cq)
    {
      _cq->Shutdown();
      void *ignored_tag;
      bool ignored_ok;
      while (_cq->Next(&ignored_tag, &ignored_ok))
      {
      }
    }
  }

  void bt_service::own_completion_queue(std::unique_ptr<grpc::ServerCompletionQueue> &&cq)
  {
    _cq = std::move(cq);
  }

  void bt_service::run()
  {
    if (!_cq)
      return;
    _pusher_manager.start();

    for (;;)
    {
      void *got_tag;
      bool ok = false;
      auto st = _cq->AsyncNext(&got_tag, &ok, gpr_time_from_millis(2000, GPR_TIMESPAN));

      if (st == CompletionQueue::SHUTDOWN)
        break;

      push_bt();

      if (st == CompletionQueue::TIMEOUT)
      {
        continue;
      }

      pusher_base::tag *tag = static_cast<pusher_base::tag *>(got_tag);

      tag->_owner->completed(tag, ok);
    }
  }

  void bt_service::push_bt()
  {
    auto n = lt::time_point::clock::now();
    if (lt::duration_cast<lt::seconds>(n - _last_push_time).count() < 2)
    {
      return;
    }
    _last_push_time = n;

    if (_ses == nullptr)
    {
      return;
    }

    _ses->post_torrent_updates();

    vector<lt::alert *> as;
    _ses->pop_alerts(&as);
    for (size_t i = 0; i < as.size(); i++)
    {

      if (auto sua = lt::alert_cast<lt::state_update_alert>(as[i]))
      {
        vector<lt::torrent_status> const &sts = sua->status;
        if (sts.size() > 0)
        {
          _pusher_manager.push_bt_status(sts);
        }
      }
      else if (auto fc = lt::alert_cast<lt::file_completed_alert>(as[i]))
      {
        _pusher_manager.push_bt_filecompleted(*fc);
      }
      else if (auto la = lt::alert_cast<lt::log_alert>(as[i]))
      {
        std::cout << la->message() << std::endl;
      }
    }
  }

  std::vector<lt::torrent_status> bt_service::get_all_bt_status()
  {
    std::vector<lt::torrent_status> ret;
    if (_ses == nullptr)
    {
      return ret;
    }
    auto ths = _ses->get_torrents();
    ret.reserve(ths.size());
    for (size_t i = 0; i < ths.size(); i++)
    {
      ret.push_back(ths[i].status());
    }
    return ret;
  }

  ::grpc::Status bt_service::InitedSession(::grpc::ServerContext *context, const ::prpc::InitedSessionReq *request, ::prpc::InitedSessionRsp *response)
  {
    (void)(context);
    (void)(request);
    response->set_inited(_ses != nullptr);
    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::InitSession(::grpc::ServerContext *context, const ::prpc::InitSessionReq *request, ::prpc::InitSessionRsp *response)
  {
    (void)(context);
    (void)(response);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    static unordered_map<string, lt::settings_pack::proxy_type_t> type_map = {
        {"socks4", lt::settings_pack::proxy_type_t::socks4},
        {"socks5", lt::settings_pack::proxy_type_t::socks5},
        {"socks5_pw", lt::settings_pack::proxy_type_t::socks5_pw},
        {"http", lt::settings_pack::proxy_type_t::http},
        {"http_pw", lt::settings_pack::proxy_type_t::http_pw},
    };

    lt::session_params sps;
    try
    {
      sps = lt::read_session_params(request->resume_data());
    }
    catch (...)
    {
      std::cout << "failed to load session resume data" << std::endl;
    }
    lt::settings_pack &sp = sps.settings;
    string const proxy_host = request->proxy_host();
    int proxy_port = request->proxy_port();
    string const proxy_type = request->proxy_type();
    int download_rate_limit = request->download_rate_limit();
    int upload_rate_limit = request->upload_rate_limit();
    int hashing_threads = request->hashing_threads();
    string listen_interfaces = request->listen_interfaces();
    if (!proxy_host.empty() && proxy_port > 0 && !proxy_type.empty() && type_map.find(proxy_type) != type_map.end())
    {
      sp.set_str(lt::settings_pack::proxy_hostname, proxy_host);
      sp.set_int(lt::settings_pack::proxy_type, type_map[proxy_type]);
      sp.set_int(lt::settings_pack::proxy_port, proxy_port);
    }
    else
    {
      sp.set_str(lt::settings_pack::proxy_hostname, "");
      sp.set_int(lt::settings_pack::proxy_type, lt::settings_pack::proxy_type_t::none);
    }
    if (listen_interfaces.size() > 0)
    {
      sp.set_str(lt::settings_pack::listen_interfaces, listen_interfaces);
    }
    sp.set_int(lt::settings_pack::download_rate_limit, download_rate_limit);
    sp.set_int(lt::settings_pack::upload_rate_limit, upload_rate_limit);
    sp.set_int(lt::settings_pack::hashing_threads, hashing_threads);
    sp.set_int(lt::settings_pack::alert_mask,
               lt::file_completed_alert::static_category | lt::log_alert::static_category);
    auto nodes = sp.get_str(lt::settings_pack::dht_bootstrap_nodes);
    if (!nodes.empty())
    {
      nodes += ", router.utorrent.com:6881";
    }
    else
    {
      nodes = "router.utorrent.com:6881";
    }
    sp.set_str(lt::settings_pack::dht_bootstrap_nodes, nodes);
    _ses = std::make_unique<lt::session>(sps);
    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::Parse(
      ::grpc::ServerContext *context, const ::prpc::DownloadRequest *request, ::prpc::DownloadRespone *response)
  {
    (void)(context);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    if (_ses == nullptr)
    {
      return ::grpc::Status(grpc::UNAVAILABLE, "");
    }
    lt::add_torrent_params params;

    switch (request->type())
    {
    case DownloadRequest_ReqType::DownloadRequest_ReqType_MagnetUri:
    {
      lt::error_code ec;

      lt::parse_magnet_uri(request->content(), params, ec);
      if (ec)
      {
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      *response->mutable_info_hash() = get_respone_info_hash(params.info_hashes);
      return ::grpc::Status::OK;
    }
    case DownloadRequest_ReqType::DownloadRequest_ReqType_Torrent:
    {
      try
      {
        string data = request->content();
        lt::load_torrent_limits cfg;
        lt::error_code ec;
        int err_pos;
        auto e = lt::bdecode(lt::span<const char>(data), ec, &err_pos, cfg.max_decode_depth, cfg.max_decode_tokens);
        params.ti = std::make_shared<lt::torrent_info>(e);
      }
      catch (std::exception const &e)
      {
        std::cout << e.what() << endl;
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      if (params.ti == nullptr)
      {
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      *response->mutable_info_hash() = get_respone_info_hash(params.ti->info_hashes());
      return ::grpc::Status::OK;
    }
    default:
    {
    }
    }
    return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
  }

  ::grpc::Status bt_service::Download(
      ::grpc::ServerContext *context, const ::prpc::DownloadRequest *request, ::prpc::DownloadRespone *response)
  {
    (void)(context);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    if (_ses == nullptr)
    {
      return ::grpc::Status(grpc::UNAVAILABLE, "");
    }
    lt::add_torrent_params params;

    switch (request->type())
    {
    case DownloadRequest_ReqType::DownloadRequest_ReqType_MagnetUri:
    {
      lt::error_code ec;

      lt::parse_magnet_uri(request->content(), params, ec);
      if (ec)
      {
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      params.save_path = request->save_path();
      break;
    }
    case DownloadRequest_ReqType::DownloadRequest_ReqType_Resume:
    {
      params = lt::read_resume_data(request->content());
      params.save_path = request->save_path();
      break;
    }
    case DownloadRequest_ReqType::DownloadRequest_ReqType_Torrent:
    {
      try
      {
        string const &data = request->content();
        lt::load_torrent_limits cfg;
        lt::error_code ec;
        int err_pos;
        auto e = lt::bdecode(lt::span<const char>(data), ec, &err_pos, cfg.max_decode_depth, cfg.max_decode_tokens);
        params.ti = std::make_shared<lt::torrent_info>(e);
        if (params.ti != nullptr)
        {
          params.info_hashes = params.ti->info_hashes();
        }
      }
      catch (std::exception const &e)
      {
        std::cout << e.what() << endl;
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      params.save_path = request->save_path();
      break;
    }
    default:
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }

    auto th = _ses->find_torrent(params.info_hashes.get_best());
    if (th.is_valid())
    {
      *response->mutable_info_hash() = get_respone_info_hash(th.info_hashes());
      if (th.flags() & lt::torrent_flags::paused)
      {
        th.resume();
      }
      return ::grpc::Status::OK;
    }

    params.extensions.push_back(torrent_plugin::build_torrent_plugin);
    auto ud = new bt_user_data{};
    ud->_ses = _ses.get();
    ud->_stop_after_got_meta = request->stop_after_got_meta();
    params.userdata = lt::client_data_t(ud);
    params.trackers.insert(params.trackers.end(), request->trackers().begin(), request->trackers().end());

    try
    {
      auto handle = _ses->add_torrent(std::move(params));
      if (handle.is_valid())
      {
        *response->mutable_info_hash() = get_respone_info_hash(handle.info_hashes());
      }
      else
      {
        return ::grpc::Status(grpc::INTERNAL, "");
      }
    }
    catch (std::exception const &)
    {
      return ::grpc::Status(grpc::INTERNAL, "");
    }

    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::RemoveTorrent(
      ::grpc::ServerContext *context, const ::prpc::RemoveTorrentReq *request, ::prpc::RemoveTorrentRes *response)
  {
    (void)(context);
    (void)(response);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    if (_ses == nullptr)
    {
      return ::grpc::Status(grpc::UNAVAILABLE, "");
    }
    auto info_hash = get_info_hash(request->info_hash());
    auto t = _ses->find_torrent(info_hash.get_best());
    if (!t.is_valid())
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    _ses->remove_torrent(t);
    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::GetMagnetUri(
      ::grpc::ServerContext *context, const ::prpc::GetMagnetUriReq *request, ::prpc::GetMagnetUriRsp *response)
  {
    (void)(context);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    if (_ses == nullptr)
    {
      return ::grpc::Status(grpc::UNAVAILABLE, "");
    }
    switch (request->type())
    {
    case GetMagnetUriReq_ReqType::GetMagnetUriReq_ReqType_Torrent:
    {
      lt::add_torrent_params params;
      try
      {
        string const &data = request->content();
        lt::load_torrent_limits cfg;
        lt::error_code ec;
        int err_pos;
        auto e = lt::bdecode(lt::span<const char>(data), ec, &err_pos, cfg.max_decode_depth, cfg.max_decode_tokens);
        params.ti = std::make_shared<lt::torrent_info>(e);
      }
      catch (std::exception const &e)
      {
        std::cout << e.what() << endl;
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      std::string uri = lt::make_magnet_uri(params);
      response->set_magnet_uri(uri);
      *response->mutable_info_hash() = get_respone_info_hash(params.ti->info_hashes());
      break;
    }
    case GetMagnetUriReq_ReqType::GetMagnetUriReq_ReqType_InfoHash:
    {
      auto th = _ses->find_torrent(get_info_hash(request->info_hash()).get_best());
      if (!th.is_valid())
      {
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "can't find torrent");
      }
      auto tf = th.torrent_file();
      if (tf == nullptr)
      {
        return ::grpc::Status(grpc::INTERNAL, "");
      }
      std::string uri = lt::make_magnet_uri(*tf);
      response->set_magnet_uri(uri);
      *response->mutable_info_hash() = get_respone_info_hash(tf->info_hashes());
      break;
    }
    default:
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::GetResumeData(::grpc::ServerContext *context, const ::prpc::GetResumeDataReq *request, ::prpc::GetResumeDataRsp *response)
  {
    (void)(context);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    if (_ses == nullptr)
    {
      return ::grpc::Status(grpc::UNAVAILABLE, "");
    }
    auto th = _ses->find_torrent(get_info_hash(request->info_hash()).get_best());
    if (!th.is_valid())
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "can't find torrent");
    }
    auto rd = th.get_resume_data(lt::torrent_handle::save_info_dict);
    auto const b = lt::write_resume_data_buf(rd);
    *response->mutable_resume_data() = std::string(b.data(), b.size());
    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::GetTorrentInfo(::grpc::ServerContext *context, const ::prpc::GetTorrentInfoReq *request, ::prpc::GetTorrentInfoRsp *response)
  {
    (void)(context);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    if (_ses == nullptr)
    {
      return ::grpc::Status(grpc::UNAVAILABLE, "");
    }
    auto th = _ses->find_torrent(get_info_hash(request->info_hash()).get_best());
    if (!th.is_valid())
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "can't find torrent");
    }
    auto ti = get_torrent_info(th);
    if (ti == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "not yet done init");
    }
    *response->mutable_torrent_info() = *ti;
    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::GetBtStatus(::grpc::ServerContext *context, const ::prpc::GetBtStatusReq *request, ::prpc::GetBtStatusRsp *response)
  {
    (void)(context);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    if (_ses == nullptr)
    {
      return ::grpc::Status(grpc::UNAVAILABLE, "");
    }
    auto th = _ses->find_torrent(get_info_hash(request->info_hash()).get_best());
    if (!th.is_valid())
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "can't find torrent");
    }
    *response->mutable_status() = get_status_respone(th.status());
    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::GetSessionParams(::grpc::ServerContext *context, const ::prpc::GetSessionParamsReq *request, ::prpc::GetSessionParamsRsp *response)
  {
    (void)(context);
    if (request == nullptr)
    {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    if (_ses == nullptr)
    {
      return ::grpc::Status(grpc::UNAVAILABLE, "");
    }
    auto sparams = _ses->session_state(lt::session_handle::save_settings | lt::session::save_dht_state | lt::session::save_extension_state | lt::session::save_ip_filter);
    auto buf = lt::write_session_params_buf(sparams);
    *response->mutable_resume_data() = std::string(buf.data(), buf.size());
    return ::grpc::Status::OK;
  }

  bt_status_pusher::bt_status_pusher(pusher_manager *owner, bt_service *ser) : pusher(owner, ser)
  {
    _ser->RequestOnBtStatus(&_context, &_stream, _ser->get_cq(), _ser->get_cq(), &_new_tag);
  }

  void bt_status_pusher::push(std::vector<lt::torrent_status> const &tss)
  {
    BtStatusRespone sr;

    auto const &req_info_hashs = _req.info_hash();
    for (auto const &ts : tss)
    {
      InfoHash info_hash = get_respone_info_hash(ts.info_hashes);
      if (info_hash.version() <= 0)
        continue;

      if (req_info_hashs.size() > 0)
      {
        auto iter = std::find_if(req_info_hashs.begin(), req_info_hashs.end(), [&info_hash](InfoHash const &i)
                                 { return i.version() == info_hash.version() && i.hash() == info_hash.hash(); });
        if (iter == req_info_hashs.end())
        {
          continue;
        }
      }
      sr.mutable_status_array()->Add(get_status_respone(ts));
    }

    write(std::move(sr));
  }

  void bt_status_pusher::completed(tag *t, bool ok)
  {
    if (t == &_new_tag && ok)
    {
      _owner->accepted_status_pusher(this);
      push(_ser->get_all_bt_status());
    }
    pusher::completed(t, ok);
  }

  void bt_status_pusher::done()
  {
    _owner->remove_status_pusher(this);
  }

  bt_filecompleted_pusher::bt_filecompleted_pusher(pusher_manager *owner, bt_service *ser) : pusher(owner, ser)
  {
    _ser->RequestOnFileCompleted(&_context, &_stream, _ser->get_cq(), _ser->get_cq(), &_new_tag);
  }

  void bt_filecompleted_pusher::push(lt::file_completed_alert const &params)
  {
    write(get_filecompleted(params));
  }

  void bt_filecompleted_pusher::completed(tag *t, bool ok)
  {
    if (t == &_new_tag && ok)
    {
      _owner->accepted_filecompleted_pusher(this);
    }
    pusher::completed(t, ok);
  }

  void bt_filecompleted_pusher::done()
  {
    _owner->remove_filecompleted_pusher(this);
  }

  pusher_manager::pusher_manager(bt_service *ser) : _ser(ser)
  {
  }

  void pusher_manager::start()
  {
    _st_pusher = std::make_unique<bt_status_pusher>(this, _ser);
    _filecompleted_pusher = std::make_unique<bt_filecompleted_pusher>(this, _ser);
  }

  void pusher_manager::push_bt_status(std::vector<lt::torrent_status> const &sts)
  {
    for (auto &pusher : _st_pushers)
    {
      pusher->push(sts);
    }
  }

  void pusher_manager::push_bt_filecompleted(lt::file_completed_alert const &params)
  {
    for (auto &pusher : _filecompleted_pushers)
    {
      pusher->push(params);
    }
  }

  void pusher_manager::accepted_status_pusher(pusher_base *pusher)
  {
    (void)(pusher);
    assert(_st_pusher != nullptr);
    _st_pushers.push_back(std::move(_st_pusher));
    _st_pusher = std::make_unique<bt_status_pusher>(this, _ser);
  }

  void pusher_manager::remove_status_pusher(pusher_base *pusher)
  {
    auto iter = std::remove_if(_st_pushers.begin(), _st_pushers.end(), [pusher](status_pusher_ptr const &p)
                               { return pusher == p.get(); });
    _st_pushers.erase(iter, _st_pushers.end());
  }

  void pusher_manager::accepted_filecompleted_pusher(pusher_base *pusher)
  {
    (void)(pusher);
    assert(_filecompleted_pusher != nullptr);
    _filecompleted_pushers.push_back(std::move(_filecompleted_pusher));
    _filecompleted_pusher = std::make_unique<bt_filecompleted_pusher>(this, _ser);
  }

  void pusher_manager::remove_filecompleted_pusher(pusher_base *pusher)
  {
    auto iter = std::remove_if(_filecompleted_pushers.begin(), _filecompleted_pushers.end(),
                               [pusher](filecompleted_pusher_ptr const &p)
                               { return pusher == p.get(); });
    _filecompleted_pushers.erase(iter, _filecompleted_pushers.end());
  }
}