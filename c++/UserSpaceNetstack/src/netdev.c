//
// Created by linlh on 2022/6/15.
//

#include <stdlib.h>
#include <string.h> //memcpy

#include "netdev.h"
#include "ethernet.h"
#include "ip.h"
#include "tuntap_if.h"

struct netdev *loop;
struct netdev *netdev;
extern int running;

static struct netdev *netdev_alloc(char *addr, char *hwaddr, uint32_t mtu) {
    struct netdev *dev = malloc(sizeof(struct netdev));

    dev->addr = ip_parse(addr);

    sscanf(hwaddr, "%hhx:%hhx:%hhx:%hhx:%hhx:%hhx", &dev->hwaddr[0],
           &dev->hwaddr[1],
           &dev->hwaddr[2],
           &dev->hwaddr[3],
           &dev->hwaddr[4],
           &dev->hwaddr[5]);
    dev->addr_len = 6;
    dev->mtu = mtu;
    return dev;
}

void netdev_init(char *addr, char *hwaddr) {
    loop = netdev_alloc("127.0.0.1", "00:00:00:00:00:00", 1500);
    netdev = netdev_alloc("10.0.0.4", "00:0c:29:6d:50:25", 1500);
}

int netdev_transmit(struct sk_buff *skb, uint8_t *dst_hw, uint16_t ethertype) {
    struct netdev *dev;
    struct eth_hdr *hdr;
    int ret = 0;

    dev = skb->dev;
    skb_push(skb, ETH_HDR_LEN);
    hdr = (struct eth_hdr *)skb->data;
    memcpy(hdr->dmac, dst_hw, dev->addr_len);
    memcpy(hdr->smac, dev->hwaddr, dev->addr_len);

    hdr->ethertype = htons(ethertype);
    eth_dbg("out", hdr);
    ret = tun_write((char *)skb->data, skb->len);
    return ret;
}

static int netdev_receive(struct sk_buff *skb) {
    struct eth_hdr *hdr = eth_hdr(skb);

    eth_dbg("in", hdr);
    switch (hdr->ethertype) {
        case ETH_P_ARP:
            break;
        case ETH_P_IP:
            break;
        case ETH_P_IPV6:
            break;
        default:
            printf("Unsupported ethertype %x\n", hdr->ethertype);
            free_skb(skb);
            break;
    }
    return 0;
}

void *netdev_rx_loop() {
    printf("inside rx loop");
    while(running) {
        struct sk_buff *skb = alloc_skb(BUFLEN);
        printf("read skb");
        if (tun_read((char *)skb->data, BUFLEN) < 0) {
            perror("ERR: Read from tun_fd");
            free_skb(skb);
            return NULL;
        }
        printf("get pack");
        netdev_receive(skb);
    }
    return NULL;
}

struct netdev *netdev_get(uint32_t sip) {
    if (netdev->addr == sip) {
        return netdev;
    } else {
        return NULL;
    }
}

void free_netdev() {
    free(loop);
    free(netdev);
}


