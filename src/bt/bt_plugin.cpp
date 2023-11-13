#include "bt_plugin.hpp"

namespace prpc
{

  torrent_plugin::torrent_plugin(lt::torrent_handle th, bt_user_data *u) : _th(th),
                                                                           _ud(u)
  {
  }
  std::shared_ptr<torrent_plugin> torrent_plugin::build_torrent_plugin(lt::torrent_handle th, lt::client_data_t u)
  {
    return std::shared_ptr<torrent_plugin>(new torrent_plugin(th, u.get<bt_user_data>()));
  }

  void torrent_plugin::on_state(lt::torrent_status::state_t st)
  {
    if (_ud->_stop_after_got_meta && _old_st == lt::torrent_status::state_t::downloading_metadata &&
        st != lt::torrent_status::state_t::downloading_metadata)
    {
      _th.pause();
    }
  }
}