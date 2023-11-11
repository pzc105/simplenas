#pragma once

#include <string>
#include <fstream>
#include <iostream>
#include "yaml-cpp/yaml.h"

inline char* getCmdOption(char** begin, char** end, const std::string& option)
{
  char** itr = std::find(begin, end, option);
  if (itr != end && ++itr != end)
  {
    return *itr;
  }
  return 0;
}

inline bool cmdOptionExists(char** begin, char** end, const std::string& option)
{
  return std::find(begin, end, option) != end;
}

namespace bt
{
  class setting
  {
  public:
    static void init(int argc, char *argv[])
    {
      m_file_name = "./bt.yml";
      if(cmdOptionExists(argv, argv + argc, "-c")) {
        m_file_name = getCmdOption(argv, argv + argc, "-c");
        std::ifstream f(m_file_name);
        if(!f.good()) {
          m_file_name = "./bt.yml";
        }
      }
      std::ifstream f(m_file_name);
      if(!f.good()) {
        m_file_name = "/etc/pnas/bt.yml";
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