cmake_minimum_required(VERSION 3.16.0 FATAL_ERROR) # Configurable policies: <= CMP0097

cmake_policy(SET CMP0091 NEW)
cmake_policy(SET CMP0092 NEW)

project(bt)

list(APPEND CMAKE_MODULE_PATH ${CMAKE_CURRENT_SOURCE_DIR}/cmake/Modules)
include(Macro)

if(CMAKE_CXX_COMPILER_ID MATCHES GNU)
  add_compile_options(
    -Wall
    -Wextra
    -Wpedantic
    -Wvla
    -Wno-noexcept-type
    -Wno-format-zero-length
    -ftemplate-depth=512
  )
elseif(MSVC)
  add_compile_options(
    /W4
    # C4251: 'identifier' : class 'type' needs to have dll-interface to be
    #        used by clients of class 'type2'
    /wd4251
    # C4268: 'identifier' : 'const' static/global data initialized
    #        with compiler generated default constructor fills the object with zeros
    /wd4268
    # C4275: non DLL-interface classkey 'identifier' used as base for
    #        DLL-interface classkey 'identifier'
    /wd4275
    # C4373: virtual function overrides, previous versions of the compiler
    #        did not override when parameters only differed by const/volatile qualifiers
    /wd4373
    # C4503: 'identifier': decorated name length exceeded, name was truncated
    /wd4503
  )
endif()

if(MSVC)
  add_definitions(-D_WIN32_WINNT=0x0601)
endif()

if(static_runtime)
  if (MSVC)
    add_compile_options(/MTd)
    set(CMAKE_MSVC_RUNTIME_LIBRARY "MultiThreaded$<$<CONFIG:Debug>:Debug>")
    message("CMAKE_MSVC_RUNTIME_LIBRARY: ${CMAKE_MSVC_RUNTIME_LIBRARY}")
    set(OPENSSL_MSVC_STATIC_RT ON)
  else()
    include(ucm_flags)
    ucm_set_runtime(STATIC)
  endif()
  set(Boost_USE_MULTITHREADED ON)
  set(Boost_USE_STATIC_LIBS ON)
  set(OPENSSL_USE_STATIC_LIBS ON)
  set(BUILD_SHARED_LIBS OFF)
endif()

set(CMAKE_CXX_STANDARD 20)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)

add_subdirectory(third_party/libtorrent)
add_subdirectory(third_party/yaml-cpp)
add_subdirectory(src/bt)