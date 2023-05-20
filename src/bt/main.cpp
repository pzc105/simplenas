#include "bt_service.hpp"
#include "grpcpp/server_builder.h"
#include "setting.hpp"

using namespace std;

int main()
{
  YAML::Node config = pset::setting::read();
  std::string server_address(config["server"]["boundAddress"].as<std::string>());

  prpc::BtService::WithAsyncMethod_OnStatus<prpc::bt_service> service;
  grpc::ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);

  service.own_completion_queue(builder.AddCompletionQueue());

  auto server = builder.BuildAndStart();

  service.run();

  cout << "end" << endl;
  return 0;
}

// char const *state(lt::torrent_status::state_t s)
// {
//   switch (s)
//   {
//   case lt::torrent_status::checking_files:
//     return "checking";
//   case lt::torrent_status::downloading_metadata:
//     return "dl metadata";
//   case lt::torrent_status::downloading:
//     return "downloading";
//   case lt::torrent_status::finished:
//     return "finished";
//   case lt::torrent_status::seeding:
//     return "seeding";
//   case lt::torrent_status::checking_resume_data:
//     return "checking resume";
//   default:
//     return "<>";
//   }
// }

// int main(/*int argc, char *argv[]*/)
// try
// {
//   string file = "test.torrent";

//   lt::settings_pack sp;
//   /*sp.set_str(lt::settings_pack::proxy_hostname, "127.0.0.1");
//   sp.set_int(lt::settings_pack::proxy_type, lt::settings_pack::socks5);
//   sp.set_int(lt::settings_pack::proxy_port, 10808);*/
//   lt::session_params sps(sp);
//   lt::session s(sp);
//   lt::add_torrent_params p = lt::parse_magnet_uri("magnet:?xt=urn:btih:19370E3FD96FB1ADA86ED5892BE5B791A2A32254&dn=John.Wick.Chapter.4.2023.HDCAM.c1nem4.x264-SUNSCREEN%5BTGx%5D&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2780%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2730%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce");
//   p.save_path = ".";
//   p.storage_mode = lt::storage_mode_allocate;
//   s.add_torrent(p);

//   while (true)
//   {
//     s.post_torrent_updates();
//     vector<lt::alert *> as;
//     s.pop_alerts(&as);
//     for (size_t i = 0; i < as.size(); i++)
//     {
//       lt::state_update_alert *sua = dynamic_cast<lt::state_update_alert *>(as[i]);
//       if (sua == nullptr)
//       {
//         continue;
//       }
//       // we only have a single torrent, so we know which one
//       // the status is for
//       if (sua->status.size() == 0)
//         continue;
//       lt::torrent_status const &s = sua->status[0];
//       std::cout << '\r' << state(s.state) << ' '
//                 << (s.download_payload_rate / 1000) << " kB/s "
//                 << (s.total_done / 1000) << " kB ("
//                 << (s.progress_ppm / 10000) << "%) downloaded ("
//                 << s.num_peers << " peers)\x1b[K";
//       std::cout.flush();
//     }
//     this_thread::sleep_for(1s);
//   }
//   return 0;
// }
// catch (std::exception const &e)
// {
//   std::cerr << "ERROR: " << e.what() << "\n";
// }
