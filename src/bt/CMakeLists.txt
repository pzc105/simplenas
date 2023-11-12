project(bt)

if(Boost_DEBUG)
  get_directory_property(_variableNames VARIABLES)
  list (SORT _variableNames)
  foreach (_variableName ${_variableNames})
    string(REGEX MATCH "BOOST.+" is_valid "${_variableName}")
    if(is_valid)
      message("${_variableName}=${${_variableName}}")
    endif()
  endforeach()
endif()
find_package(Boost COMPONENTS REQUIRED)

include_directories(./)
include_directories(./prpc)
include_directories(./prpc/google/api)

add_executable(bt
    setting.cpp
    translate.cpp
    bt_service.cpp
    bt_plugin.cpp
    main.cpp)

set(src-pattern "prpc/*.cc" "prpc/google/api/*.cc")
file(GLOB prpc-cpp-sources ${src-pattern})
target_sources(bt PRIVATE
    ${prpc-cpp-sources})

find_package(OpenSSL)
if(TARGET OpenSSL::SSL)
  get_target_property(SSL_LIB OpenSSL::SSL IMPORTED_LOCATION)
  message("SSL library path: ${SSL_LIB}")
  get_target_property(SSL_INCLUDE_DIRS OpenSSL::SSL INTERFACE_INCLUDE_DIRECTORIES)
  message("SSL header files path: ${SSL_INCLUDE_DIRS}")
else()
  message("not found openssl")
endif()

find_package(protobuf CONFIG REQUIRED)
find_package(gRPC CONFIG REQUIRED)

target_link_libraries(bt PRIVATE torrent-rasterbar yaml-cpp ${Boost_LIBRARIES})

target_link_libraries(bt PRIVATE gRPC::grpc++_reflection protobuf::libprotobuf)

target_compile_features(bt
  PUBLIC
    cxx_std_14
    cxx_attribute_deprecated
    cxx_binary_literals
    cxx_contextual_conversions
    cxx_decltype_auto
    cxx_digit_separators
    cxx_generic_lambdas
    cxx_lambda_init_captures
    cxx_relaxed_constexpr
    cxx_variable_templates
)

set_property(TARGET bt PROPERTY
             MSVC_RUNTIME_LIBRARY "MultiThreaded$<$<CONFIG:Debug>:Debug>")