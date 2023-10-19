#pragma once

#include <string>
#include <fstream>
#include <iostream>
#include <boost/program_options.hpp>
#include "yaml-cpp/yaml.h"

namespace bt
{
  class setting
  {
  public:
    static YAML::Node read(int argc, char *argv[])
    {
      boost::program_options::options_description opts("all options");
      boost::program_options::variables_map vm;
      opts.add_options()("c", boost::program_options::value<std::string>(), "config path");
      try
      {
        boost::program_options::store(boost::program_options::parse_command_line(argc, argv, opts), vm);
      }
      catch (...)
      {
        std::cout << "输入参数有问题" << std::endl;
        exit(-1);
      }

      std::string file_name = "./bt.yml";
      if(vm.count("c"))
      {
        file_name = vm["c"].as<std::string>();
      }
      return YAML::LoadFile(file_name);
    }
  };
}