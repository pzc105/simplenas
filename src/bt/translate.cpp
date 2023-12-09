#include "translate.hpp"
#include "libtorrent/write_resume_data.hpp"
#include <memory>

namespace prpc
{
  lt::info_hash_t get_info_hash(InfoHash const &i)
  {
    if (i.version() == 1 && static_cast<std::size_t>(lt::sha1_hash::size()) == i.hash().size())
    {
      return lt::info_hash_t(lt::sha1_hash(i.hash().data()));
    }
    if (i.version() == 2 && static_cast<std::size_t>(lt::sha256_hash::size()) == i.hash().size())
    {
      return lt::info_hash_t(lt::sha256_hash(i.hash().data()));
    }
    return lt::info_hash_t();
  }
  InfoHash get_respone_info_hash(lt::info_hash_t const &i)
  {
    InfoHash info_hash;
    if (i.has_v2())
    {
      info_hash.set_version(2);
      info_hash.set_hash(i.get(lt::protocol_version::V2).to_string());
    }
    else if (i.has_v1())
    {
      info_hash.set_version(1);
      info_hash.set_hash(i.get(lt::protocol_version::V1).to_string());
    }
    return info_hash;
  }

  TorrentStatus get_status_respone(lt::torrent_status const &ts)
  {
    TorrentStatus ret;
    *ret.mutable_info_hash() = get_respone_info_hash(ts.info_hashes);
    ret.set_name(ts.name);
    ret.set_download_payload_rate(ts.download_payload_rate);
    ret.set_total_done(ts.total_done);
    ret.set_total(ts.total);
    ret.set_progress(ts.progress);
    ret.set_num_peers(ts.num_peers);
    ret.set_state(static_cast<prpc::BtStateEnum>(ts.state));
    return ret;
  }

  std::unique_ptr<TorrentInfo> get_torrent_info(lt::torrent_handle const &th)
  {
    auto tf = th.torrent_file();
    if (tf == nullptr)
    {
      return nullptr;
    }
    auto ret = std::make_unique<TorrentInfo>();
    *ret->mutable_info_hash() = get_respone_info_hash(th.info_hashes());
    ret->set_name(tf->name());
    auto st = th.status(lt::torrent_handle::query_save_path);
    ret->set_save_path(st.save_path);

    auto const &storage = tf->files();
    auto range = storage.file_range();
    for (auto const i : range)
    {
      auto file = ret->add_files();
      file->set_index(i);
      file->set_name(storage.file_path(i));
      file->set_total_size(storage.file_size(i));
    }
    ret->set_total_size(tf->total_size());
    ret->set_piece_length(tf->piece_length());
    ret->set_num_pieces(tf->num_pieces());
    return ret;
  }

  FileCompletedRes get_filecompleted(lt::file_completed_alert const &params)
  {
    FileCompletedRes res;
    auto const &th = params.handle;
    *res.mutable_info_hash() = get_respone_info_hash(th.info_hashes());
    res.set_file_index(params.index);
    return res;
  }

  static std::string print_endpoint(lt::tcp::endpoint const &ep)
  {
    using namespace lt;
    char buf[200];
    address const &addr = ep.address();
    if (addr.is_v6())
      std::snprintf(buf, sizeof(buf), "[%s]:%d", addr.to_string().c_str(), ep.port());
    else
      std::snprintf(buf, sizeof(buf), "%s:%d", addr.to_string().c_str(), ep.port());
    return buf;
  }

  PeerInfo get_rsp(lt::peer_info const &p)
  {
    PeerInfo res;
    res.set_client(p.client);
    res.set_total_download(p.total_download);
    res.set_total_upload(p.total_upload);
    res.set_flags(p.flags);
    res.set_source(p.source);
    res.set_up_speed(p.up_speed);
    res.set_down_speed(p.down_speed);
    res.set_payload_up_speed(p.payload_up_speed);
    res.set_payload_down_speed(p.payload_down_speed);
    res.set_pid(p.pid.to_string());
    res.set_queue_bytes(p.queue_bytes);
    res.set_connection_type(p.connection_type);
    res.set_download_rate_peak(p.download_rate_peak);
    res.set_upload_rate_peak(p.upload_rate_peak);
    res.set_peer_addr(print_endpoint(p.ip));
    res.set_num_pieces(p.num_pieces);
    return res;
  }
}