#include <unordered_map>
#include <string_view>
#include "bt_service.hpp"
#include "libtorrent/read_resume_data.hpp"
#include "translate.hpp"

using namespace std;
using namespace grpc;

namespace prpc
{
  bt_service::bt_service() : _pusher_manager(this)
  {
    YAML::Node config = pset::setting::read();
    YAML::Node bt_config = config["bt"];

    static unordered_map<string, lt::settings_pack::proxy_type_t> type_map = {
        {"socks4", lt::settings_pack::proxy_type_t::socks4},
        {"socks5", lt::settings_pack::proxy_type_t::socks5},
        {"socks5_pw", lt::settings_pack::proxy_type_t::socks5_pw},
        {"http", lt::settings_pack::proxy_type_t::http},
        {"http_pw", lt::settings_pack::proxy_type_t::http_pw},
    };

    lt::settings_pack sp;
    string const proxy_host = bt_config["proxy_hostname"].as<string>();
    int proxy_port = std::atoi(bt_config["proxy_port"].as<string>().c_str());
    string const proxy_type = bt_config["proxy_type"].as<string>();
    int download_rate_limit = std::atoi(bt_config["download_rate_limit"].as<string>().c_str());
    int upload_rate_limit = std::atoi(bt_config["upload_rate_limit"].as<string>().c_str());
    int hashing_threads = std::atoi(bt_config["hashing_threads"].as<string>().c_str());
    if (!proxy_host.empty() && proxy_port > 0 && !proxy_type.empty() && type_map.find(proxy_type) != type_map.end()) {
      sp.set_str(lt::settings_pack::proxy_hostname, proxy_host);
      sp.set_int(lt::settings_pack::proxy_type, type_map[proxy_type]);
      sp.set_int(lt::settings_pack::proxy_port, proxy_port);
    }
    sp.set_int(lt::settings_pack::download_rate_limit, download_rate_limit);
    sp.set_int(lt::settings_pack::upload_rate_limit, upload_rate_limit);
    sp.set_int(lt::settings_pack::hashing_threads, hashing_threads);
    sp.set_int(lt::settings_pack::alert_mask, lt::file_completed_alert::static_category);
    lt::session_params sps(sp);
    _ses = std::make_unique<lt::session>(sps);
  }

  bt_service::~bt_service()
  {
    if (_cq) {
      _cq->Shutdown();
      void *ignored_tag;
      bool ignored_ok;
      while (_cq->Next(&ignored_tag, &ignored_ok)) { }
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

    for (;;) {
      void *got_tag;
      bool ok = false;
      auto st = _cq->AsyncNext(&got_tag, &ok, gpr_time_from_millis(2000, GPR_TIMESPAN));

      if (st == CompletionQueue::SHUTDOWN)
        break;

      push_bt();
      
      if (st == CompletionQueue::TIMEOUT) {
        continue;
      }

      pusher_base::tag* tag = static_cast<pusher_base::tag*>(got_tag);

      tag->_owner->completed(tag, ok);
    }
  }

  void bt_service::push_bt()
  {
    auto n = lt::time_point::clock::now();
    if (lt::duration_cast<lt::seconds>(n - _last_push_time).count() < 2) {
      return;
    }

    _last_push_time = n;

    _ses->post_torrent_updates();
    auto tss = _ses->get_torrents();
    for (auto const& t : tss) {
      t.save_resume_data(lt::torrent_handle::only_if_modified | lt::torrent_handle::save_info_dict);
    }

    vector<lt::alert *> as;
    _ses->pop_alerts(&as);
    for (size_t i = 0; i < as.size(); i++) {

      if (auto sua = lt::alert_cast<lt::state_update_alert>(as[i])) {
        vector<lt::torrent_status> const& sts = sua->status;
        if (sts.size() > 0) {
          _pusher_manager.push_bt_status(sts);
        }
      }
      else if (auto srd = lt::alert_cast<lt::save_resume_data_alert>(as[i])) {
        _pusher_manager.push_bt_infos(srd->params);
      }
      else if (auto fc = lt::alert_cast<lt::file_completed_alert>(as[i])) {
        _pusher_manager.push_bt_filecompleted(*fc);
      }
    }
  }

  ::grpc::Status bt_service::Parse(::grpc::ServerContext* context, const::prpc::DownloadRequest* request, ::prpc::DownloadRespone* response)
  {
    (void)(context);
    lt::add_torrent_params params;

    switch (request->type())
    {
    case DownloadRequest_ReqType::DownloadRequest_ReqType_MagnetUri: {
      lt::error_code ec;

      lt::parse_magnet_uri(request->content(), params, ec);
      if (ec) {
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      *response->mutable_info_hash() = get_respone_info_hash(params.info_hashes);
      return ::grpc::Status::OK;
    }
    case DownloadRequest_ReqType::DownloadRequest_ReqType_Torrent: {
      try {
        string data = request->content();
        lt::load_torrent_limits cfg;
        lt::error_code ec;
        int err_pos;
        auto e = lt::bdecode(lt::span<const char>(data), ec, &err_pos
          , cfg.max_decode_depth, cfg.max_decode_tokens);
        params.ti = std::make_shared<lt::torrent_info>(e);
      }
      catch (std::exception const& e) {
        std::cout << e.what() << endl;
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      *response->mutable_info_hash() = get_respone_info_hash(params.info_hashes);
      return ::grpc::Status::OK;
    }
    default: {

    }
    }
    return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
  }

  ::grpc::Status bt_service::Download(::grpc::ServerContext* context, const ::prpc::DownloadRequest* request, ::prpc::DownloadRespone* response)
  {
    (void)(context);
    lt::add_torrent_params params;

    switch (request->type())
    {
    case DownloadRequest_ReqType::DownloadRequest_ReqType_MagnetUri: {
      lt::error_code ec;

      lt::parse_magnet_uri(request->content(), params, ec);
      if (ec) {
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      params.save_path = request->save_path();
      break;
    }
    case DownloadRequest_ReqType::DownloadRequest_ReqType_Resume:{
      params = lt::read_resume_data(request->content());
      params.save_path = request->save_path();
      break;
    }
    case DownloadRequest_ReqType::DownloadRequest_ReqType_Torrent: {
      try {
        string data = request->content();
        lt::load_torrent_limits cfg;
        lt::error_code ec;
        int err_pos;
        auto e = lt::bdecode(lt::span<const char>(data), ec, &err_pos
          , cfg.max_decode_depth, cfg.max_decode_tokens);
        params.ti = std::make_shared<lt::torrent_info>(e);
      }
      catch (std::exception const& e) {
        std::cout << e.what() << endl;
        return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
      }
      params.save_path = request->save_path();
      break;
    }
    default:
     return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }

    try {
      auto handle = _ses->add_torrent(std::move(params));
      *response->mutable_info_hash() = get_respone_info_hash(handle.info_hashes());
    }
    catch (std::exception const&) {
      return ::grpc::Status(grpc::INTERNAL, "");
    }

    return ::grpc::Status::OK;
  }

  ::grpc::Status bt_service::RemoveTorrent(
    ::grpc::ServerContext* context
    , const::prpc::RemoveTorrentReq* request
    , ::prpc::RemoveTorrentRes* response)
  {
    (void)(context);
    (void)(response);
    if (request == nullptr) {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    auto info_hash = get_info_hash(request->info_hash());
    auto t = _ses->find_torrent(info_hash.get_best());
    if (!t.is_valid()) {
      return ::grpc::Status(grpc::INVALID_ARGUMENT, "");
    }
    _ses->remove_torrent(t);
    return ::grpc::Status::OK;
  }

  bt_status_pusher::bt_status_pusher(pusher_manager *owner, bt_service *ser) : pusher(owner, ser)
  {
    _ser->RequestOnStatus(&_context, &_stream, _ser->get_cq(), _ser->get_cq(), &_new_tag);
  }

  void bt_status_pusher::push(std::vector<lt::torrent_status> const &tss)
  {
    StatusRespone sr;

    auto const &req_info_hashs = _req.info_hash();
    for (auto const &ts : tss) {
      InfoHash info_hash = get_respone_info_hash(ts.info_hashes);
      if (info_hash.version() <= 0)
        continue;

      if (req_info_hashs.size() > 0) {
        auto iter = std::find_if(req_info_hashs.begin(), req_info_hashs.end(), [&info_hash](InfoHash const &i)
                                 { return i.version() == info_hash.version() && i.hash() == info_hash.hash(); });
        if (iter == req_info_hashs.end()) {
          continue;
        }
      }
      sr.mutable_status_array()->Add(get_status_respone(ts));
    }

    write(std::move(sr));
  }

  void bt_status_pusher::completed(tag* t, bool ok)
  {
    if (t == &_new_tag && ok) {
      _owner->accepted_status_pusher(this);
    }
    pusher::completed(t, ok);
  }

  void bt_status_pusher::done()
  {
    _owner->remove_status_pusher(this);
  }

  bt_info_pusher::bt_info_pusher(pusher_manager* owner, bt_service* ser) : pusher(owner, ser)
  {
    _ser->RequestOnTorrentInfo(&_context, &_stream, _ser->get_cq(), _ser->get_cq(), &_new_tag);
  }

  void bt_info_pusher::push(lt::add_torrent_params const& params)
  {
    TorrentInfoRes tr;
    *tr.mutable_ti() = get_torrent_info(params);
    write(std::move(tr));
  }

  void bt_info_pusher::completed(tag* t, bool ok)
  {
    if (t == &_new_tag && ok) {
      _owner->accepted_btinfo_pusher(this);
    }
    pusher::completed(t, ok);
  }

  void bt_info_pusher::done()
  {
    _owner->remove_btinfo_pusher(this);
  }


  bt_filecompleted_pusher::bt_filecompleted_pusher(pusher_manager* owner, bt_service* ser) : pusher(owner, ser)
  {
    _ser->RequestOnFileCompleted(&_context, &_stream, _ser->get_cq(), _ser->get_cq(), &_new_tag);
  }

  void bt_filecompleted_pusher::push(lt::file_completed_alert const& params)
  {
    write(get_filecompleted(params));
  }

  void bt_filecompleted_pusher::completed(tag* t, bool ok)
  {
    if (t == &_new_tag && ok) {
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
    _info_pusher = std::make_unique<bt_info_pusher>(this, _ser);
    _filecompleted_pusher = std::make_unique<bt_filecompleted_pusher>(this, _ser);
  }

  void pusher_manager::push_bt_status(std::vector<lt::torrent_status> const &sts)
  {
    for (auto &pusher : _st_pushers) {
      pusher->push(sts);
    }
  }

  void pusher_manager::push_bt_infos(lt::add_torrent_params const& params)
  {
    for (auto& pusher : _info_pushers) {
      pusher->push(params);
    }
  }

  void pusher_manager::push_bt_filecompleted(lt::file_completed_alert const& params)
  {
    for (auto& pusher : _filecompleted_pushers) {
      pusher->push(params);
    }
  }

  void pusher_manager::accepted_status_pusher(pusher_base* pusher)
  {
    (void)(pusher);
    assert(_st_pusher != nullptr);
    _st_pushers.push_back(std::move(_st_pusher));
    _st_pusher = std::make_unique<bt_status_pusher>(this, _ser);
  }

  void pusher_manager::remove_status_pusher(pusher_base* pusher)
  {
    auto iter = std::remove_if(_st_pushers.begin(), _st_pushers.end(), [pusher](status_pusher_ptr const& p)
      { return pusher == p.get(); });
    _st_pushers.erase(iter, _st_pushers.end());
  }

  void pusher_manager::accepted_btinfo_pusher(pusher_base* pusher)
  {
    (void)(pusher);
    assert(_info_pusher != nullptr);
    _info_pushers.push_back(std::move(_info_pusher));
    _info_pusher = std::make_unique<bt_info_pusher>(this, _ser);
  }

  void pusher_manager::remove_btinfo_pusher(pusher_base *pusher)
  {
    auto iter = std::remove_if(_info_pushers.begin(), _info_pushers.end(), [pusher](btinfo_pusher_ptr const& p)
      { return pusher == p.get(); });
    _info_pushers.erase(iter, _info_pushers.end());
  }

  void pusher_manager::accepted_filecompleted_pusher(pusher_base* pusher)
  {
    (void)(pusher);
    assert(_filecompleted_pusher != nullptr);
    _filecompleted_pushers.push_back(std::move(_filecompleted_pusher));
    _filecompleted_pusher = std::make_unique<bt_filecompleted_pusher>(this, _ser);
  }

  void pusher_manager::remove_filecompleted_pusher(pusher_base* pusher)
  {
    auto iter = std::remove_if(_filecompleted_pushers.begin(), _filecompleted_pushers.end(),
      [pusher](filecompleted_pusher_ptr const& p) { return pusher == p.get(); });
    _filecompleted_pushers.erase(iter, _filecompleted_pushers.end());
  }
}