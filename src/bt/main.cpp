#include "bt_service.hpp"
#include "grpcpp/server_builder.h"
#include "setting.hpp"

using namespace std;

int main(int argc, char *argv[])
{
  bt::setting::init(argc, argv);
  YAML::Node config = bt::setting::read();
  std::string server_address(config["server"]["boundAddress"].as<std::string>());

  prpc::bt_service service;
  grpc::ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);

  service.own_completion_queue(builder.AddCompletionQueue());

  auto server = builder.BuildAndStart();

  service.run();

  cout << "end" << endl;
  return 0;
}