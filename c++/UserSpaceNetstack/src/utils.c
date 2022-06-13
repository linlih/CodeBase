//
// Created by linlh on 2022/6/13.
//

#include <stdio.h> // va_list
#include <stdarg.h> // va_start
#include <stdlib.h> // system

#include "utils.h"

int run_cmd(char *cmd, ...) {
    va_list ap;
    char buf[CMDBUFLEN];
    va_start(ap, cmd);
    vsnprintf(buf, CMDBUFLEN, cmd, ap);
    va_end(ap);

    printf("EXEC: %s\n", buf);

    return system(buf);
}
