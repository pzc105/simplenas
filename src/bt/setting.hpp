#pragma once

#include <string>
#include <fstream>
#include "yaml-cpp/yaml.h"

namespace pset
{
  class setting
  {
  public:
    static YAML::Node read()
    {
      std::string file_name = "local_setting.yml";
      std::ifstream fin(file_name);
      if(!fin.is_open()) {
        file_name = "setting.yml";
      }
      return YAML::LoadFile(file_name);
    }
  };
}