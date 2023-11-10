#include "prpc/bt.pb.h"
#include "libtorrent/torrent_info.hpp"
#include "libtorrent/torrent_status.hpp"
#include "libtorrent/add_torrent_params.hpp"
#include "libtorrent/alert_types.hpp"

namespace prpc
{
  lt::info_hash_t get_info_hash(InfoHash const& i);
  InfoHash get_respone_info_hash(lt::info_hash_t const& i);
  TorrentStatus get_status_respone(lt::torrent_status const& ts);
  TorrentInfo get_torrent_info(lt::torrent_handle const& th);
  FileCompletedRes get_filecompleted(lt::file_completed_alert const& params);


}