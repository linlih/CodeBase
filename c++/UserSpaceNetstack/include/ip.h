//
// Created by linlh on 2022/6/15.
//

#ifndef NETSTACK_IP_H
#define NETSTACK_IP_H

#include <stdint.h>
#include <arpa/inet.h> // inet_pton
#include <stdlib.h> // exit

#include "skbuff.h"
#include "ethernet.h"

#define IPV4 0x04
#define IP_TCP 0x06
#define ICMPV4 0x01

#define IP_HDR_LEN sizeof(struct iphdr)
#define ip_len(ip) (ip->len - (ip->ihl * 4))

#if DEBUG_IP
#define ip_dbg(msg, hdr) \
    do{                  \
    print_debug("ip "msg" (ihl")\
    }
#else
#define ip_dbg(msg, hdr)
#endif

struct iphdr {
    uint8_t ihl : 4;
    uint8_t version : 4;
    uint8_t tos;
    uint16_t len;
    uint16_t id;
    uint16_t frag_offset;
    uint8_t ttl;
    uint8_t proto;
    uint16_t csum;
    uint32_t saddr;
    uint32_t daddr;
    uint8_t data[];
} __attribute__((packed));

static inline struct iphdr *ip_hdr(const struct sk_buff *skb) {
    return (struct iphdr*)(skb->head + ETH_HDR_LEN);
}

static inline uint32_t ip_parse(char *addr) {
    uint32_t dst = 0;
    if (inet_pton(AF_INET, addr, &dst) != 1) {
        perror("ERR: Parsing inet address failed");
        exit(1);
    }
    return ntohl(dst);
}

#endif //NETSTACK_IP_H
