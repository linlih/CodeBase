//
// Created by linlh on 2022/6/11.
//
#include <pthread.h>

#include "tuntap_if.h"
#include "netdev.h"
#include "route.h"
#include "utils.h"

int running = 1; // 全局启动标志

#define THREAD_CORE   0
#define THREAD_TIMERS 1
#define THREAD_IPC    2
#define THREAD_SIGNAL 3
static pthread_t threads[4];

static void create_thread(pthread_t id, void *(*func) (void *)) {
    if (pthread_create(&threads[id], NULL, func, NULL) != 0) {
        print_err("Could not create core thread\n");
    }
}

static run_threads() {
    create_thread(THREAD_CORE, netdev_rx_loop);
}

int main() {
    tun_init();
    netdev_init();
    route_init();
    run_threads();
    while (1);
}
