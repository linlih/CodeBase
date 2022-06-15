//
// Created by linlh on 2022/6/11.
//

#include <stdlib.h> // calloc
#include <fcntl.h> // open
#include <stdio.h> // perror
#include <net/if.h> // ifreq struct
#include <sys/ioctl.h>
#include <string.h> // strncpy
#include <unistd.h> // close
#include <linux/if_tun.h> // IFF_TAP, IFF_NO_PI, TUNSETIFF

#include "utils.h"

static int tun_alloc(char *dev) {
    // Linux 支持一些标准的 ioctl 来配置网络设备，他们可以在任何 Socket 文件描述符上使用，大部分使用的是 ifreq 结构体来配置
    struct ifreq ifr;
    int fd, err;
    if ((fd = open("/dev/net/tap", O_RDWR)) < 0) {
        perror("Cannot open TUN/TAP device\n"
               "Make sure one exists with "
               "'$ mknod /dev/net/tap c 10 200'");
        exit(1);
    }
    memset(&ifr, 0, sizeof(ifr));

    ifr.ifr_flags = IFF_TAP | IFF_NO_PI;
    if (*dev) {
        strncpy(ifr.ifr_name, dev, IFNAMSIZ);
    }

    if ((err = ioctl(fd, TUNSETIFF, (void*)&ifr)) < 0) {
        perror("ERR: Could not ioctl on");
        close(fd);
        exit(err);
    }
    strcpy(dev, ifr.ifr_name);
    return fd;
}

static int set_if_route(char* dev, char* cidr) {
    return run_cmd("ip route add dev %s %s", dev, cidr);
}

static int set_if_address(char* dev, char* cidr) {
    return run_cmd("ip address add dev %s local %s", dev, cidr);
}

static int set_if_up(char *dev) {
    return run_cmd("ip link set dev %s up", dev);
}


static int tun_fd;
static  char* dev;

char* tapaddr = "10.0.0.5";
char* taproute = "10.0.0.0/24";

void tun_init() {
    dev = calloc(10, 1);
    tun_fd = tun_alloc(dev);
    if (set_if_up(dev) != 0) {
        print_err("ERROR when setting up if\n");
    }
    if (set_if_route(dev, taproute) != 0) {
        print_err("ERROR when setting route for if\n");
    }
    if (set_if_address(dev, tapaddr) != 0) {
        print_err("ERROR when setting addr for if\n");
    }
}

int tun_write(char* buf, int len) {
    return write(tun_fd, buf, len);
}

int tun_read(char* buf, int len) {
    return read(tun_fd, buf, len);
}

void free_tun() {
    free(dev);
}