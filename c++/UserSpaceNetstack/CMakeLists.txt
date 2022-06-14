cmake_minimum_required(VERSION 3.21)
project(netstack)

set(CMAKE_CXX_STANDARD 14)

include_directories(include)

add_executable(netstack
        src/main.c
        src/tuntap_if.c src/utils.c include/utils.h include/tuntap_if.h include/ethernet.h include/skbuff.h include/netdev.h include/list.h include/route.h)
