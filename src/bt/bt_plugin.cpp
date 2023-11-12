#include "bt_plugin.hpp"

namespace prpc
{

  torrent_plugin::torrent_plugin(lt::torrent_handle const &th, lt::session *ses) : _th(th),
                                                                                   _ses(ses)
  {
  }
  std::shared_ptr<torrent_plugin> torrent_plugin::build_torrent_plugin(lt::torrent_handle const &th, lt::client_data_t u)
  {
    return std::shared_ptr<torrent_plugin>(new torrent_plugin(th, u.get<lt::session>()));
  }


  void torrent_plugin::on_state(lt::torrent_status::state_t)
  {
    
  }
}