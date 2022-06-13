//
// Created by linlh on 2022/6/13.
//

#ifndef NETSTACK_UTILS_H
#define NETSTACK_UTILS_H

#define CMDBUFLEN 100

#define print_err(str, ...) \
    fprintf(stderr, str, ##__VA_ARGS__);

int run_cmd(char* cmd, ...);

#endif //NETSTACK_UTILS_H
