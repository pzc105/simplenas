#include "libtorrent/torrent_handle.hpp"
#include "libtorrent/extensions.hpp"
#include <memory>

namespace prpc
{
  class torrent_plugin : public lt::torrent_plugin
  {
  public:
    static std::shared_ptr<torrent_plugin> build_torrent_plugin(lt::torrent_handle const &th, lt::client_data_t);

  private:
    torrent_plugin(lt::torrent_handle const& th, lt::session *ses);
    void on_state(lt::torrent_status::state_t) override;

  private:
    lt::torrent_handle const& _th;
    lt::session *_ses;
  };
}