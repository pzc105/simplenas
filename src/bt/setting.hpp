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
    static void init(int argc, char *argv[])
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

      m_file_name = "./bt.yml";
      if(vm.count("c"))
      {
        m_file_name = vm["c"].as<std::string>();
      }
    }

    static YAML::Node read()
    {
      return YAML::LoadFile(m_file_name);
    }
  private:
    static std::string m_file_name;
  };
}