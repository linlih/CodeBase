//
// Created by linlh on 2022/6/13.
//

#ifndef NETSTACK_TUNTAP_IF_H
#define NETSTACK_TUNTAP_IF_H

void tun_init();
int tun_read(char* buf, int len);
int tun_write(char* buf, int len);
void free_tun();

#endif //NETSTACK_TUNTAP_IF_H
