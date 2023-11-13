#include "libtorrent/torrent_handle.hpp"
#include "libtorrent/extensions.hpp"
#include <memory>

namespace prpc
{

  struct bt_user_data
  {
    lt::session *_ses;
    bool _stop_after_got_meta;
  };

  class torrent_plugin : public lt::torrent_plugin
  {
  public:
    static std::shared_ptr<torrent_plugin> build_torrent_plugin(lt::torrent_handle th, lt::client_data_t);

  private:
    torrent_plugin(lt::torrent_handle th, bt_user_data *ses);
    void on_state(lt::torrent_status::state_t) override;

  private:
    lt::torrent_handle _th;
    std::unique_ptr<bt_user_data> _ud;
    lt::torrent_status::state_t _old_st;
  };
}