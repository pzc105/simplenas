#include "translate.hpp"
#include "libtorrent/write_resume_data.hpp"

namespace prpc
{
  lt::info_hash_t get_info_hash(InfoHash const& i)
  {
    if (i.version() == 1 && static_cast<std::size_t>(lt::sha1_hash::size()) == i.hash().size()) {
      return lt::info_hash_t(lt::sha1_hash(i.hash().data()));
    }
    if (i.version() == 2 && static_cast<std::size_t>(lt::sha256_hash::size()) == i.hash().size()) {
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

  TorrentInfo get_torrent_info(lt::add_torrent_params const &params)
  {
    TorrentInfo ret;
    *ret.mutable_info_hash() = get_respone_info_hash(params.info_hashes);
    ret.set_name(params.name);
    ret.set_save_path(params.save_path);
    auto const& ti = params.ti;
    if (ti && ti->is_valid())
    {
      auto const& storage = ti->files();
      int32_t index = 0;
      for (auto const &f : storage)
      {
        auto file = ret.add_files();
        file->set_index(index++);
        file->set_name(storage.file_path(f));
        file->set_total_size(storage.file_size(f));
      }
      ret.set_total_size(ti->total_size());
      ret.set_piece_length(ti->piece_length());
      ret.set_num_pieces(ti->num_pieces());
    }
    ret.set_state(static_cast<prpc::BtStateEnum>(params.state));

    auto const b = lt::write_resume_data_buf(params);
    *ret.mutable_resume_data() = std::string(b.data(), b.size());
    return ret;
  }

  FileCompletedRes get_filecompleted(lt::file_completed_alert const& params)
  {
    FileCompletedRes res;
    auto const& th = params.handle;
    *res.mutable_info_hash() = get_respone_info_hash(th.info_hashes());
    res.set_file_index(params.index);
    return res;
  }
}