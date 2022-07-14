//
// Created by linlh on 22-6-25.
//
#include <stdio.h>
#include <malloc.h>


// 内存实验
// 当分配的内存大小是小于128KB的时候，malloc使用的是brk从堆上直接分配的
//
int main() {
    printf("使用:cat /proc/%d/maps 查看进程的内存分配\n", getpid());
    // 分配一个字节
    void *addr = malloc(129*1024);
    printf("分配一个字节的内存起始地址: %p\n", addr);
    printf("使用:cat /proc/%d/maps 查看进程分配1字节后的内存分配\n", getpid());

    getchar();
    // 释放内存
    free(addr);
    printf("使用:cat /proc/%d/maps 查看进程释放1字节后的内存分配\n", getpid());

    getchar();
    return 0;
}