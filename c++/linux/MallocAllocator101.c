#include <pthread.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

/*
  一个程序的内存结构：

  Address
  high ┌-----------------┐
       │      stack      │
       │-----------------│
       │        ↓        │
       │                 │
       │        ↑        │
       ├-----------------┤← program break，brk
       │       heap      │
       │-----------------│
       │       bss       │ uninitialized variables
       │-----------------│
       │       data      │ initialized variables
       │-----------------│
       │       text      │ instruction
    0  └-----------------┘
 */

typedef char ALIGN[16];

union header {
  // 这个struct在32位上12字节，在64位是16字节
  struct {
    size_t size;
    unsigned is_free;
    union header *next;
  } s;
  // 强制让 header 头部对齐到 16 字节
  ALIGN stub;
};

typedef union header header_t;

header_t *head = NULL, *tail = NULL;

pthread_mutex_t global_malloc_lock;

// 从链表头开始寻找一个满足大小的块
header_t *get_free_block(size_t size) {
  header_t *curr = head;
  while (curr) {
    if (curr->s.is_free && curr->s.size >= size) {
      return curr;
    }
    curr = curr->s.next;
  }
  return NULL;
}

void free(void *block) {
  header_t *header, *tmp;
  // 记录的是程序数据段结尾的位置
  void *programbreak;
  if (!block) {
    return;
  }
  pthread_mutex_lock(&global_malloc_lock);
  header = (header_t *)block - 1;
  // sbrk(0) 返回当前程序的 break 地址
  programbreak = sbrk(0);

  // 如果当前的块刚好的brk的下面一块的时候，就将其归还给操作系统
  // 其他场景下是标记为 free，挂在链上
  if ((char *)block + header->s.size == programbreak) {
    if (head == tail) {
      head = tail = NULL;
    } else {
      tmp = head;
      while (tmp) {
        if (tmp->s.next == tail) {
          tmp->s.next = NULL;
          tail = tmp;
        }
        tmp = tmp->s.next;
      }
    }
    /*
     * sbrk 传入一个非 0 的地址，减小程序的 break 地址
     * 这样内存就会归还给操作系统了
     * */
    sbrk(0 - header->s.size - sizeof(header_t));
    pthread_mutex_unlock(&global_malloc_lock);
    return;
  }
  header->s.is_free = 1;
  pthread_mutex_unlock(&global_malloc_lock);
}

// malloc 中不能调用 printf，因为 printf 函数本身会调用 malloc，这样就会导致死锁
// 可行的方式就是修改下 malloc 的名字，如下面修改为 Malloc
// 如果是使用小写的 malloc，会看到分配会多分配一个 1024 的块，这个场景下只有在调用了 printf 会出现
// 原因还是一样的，printf 会调用 malloc 函数，如果我们自己定义的函数是 malloc，那么也会被 printf 调用到
// printf 调用就会分配一块 1024 的内存空间，这个内存块也会交给我们自己定义的 allocator 来管理
void *Malloc(size_t size) {
  size_t total_size;
  void *block;
  header_t *header;
  if (!size) {
    return NULL;
  }
  pthread_mutex_lock(&global_malloc_lock);
  header = get_free_block(size);
  if (header) {
    // 如果找到空闲的块，那么可以直接返回
    header->s.is_free = 0;
    pthread_mutex_unlock(&global_malloc_lock);
    // 这个地方+1，是加上了header大小，也就是16字节，这个需要注意
    return (void *)(header + 1);
  }
  total_size = sizeof(header_t) + size;
  //printf("malloc from os: %ld\n", total_size);
  // 找不到满足条件的内存，则调用 sbrk 来申请内存
  block = sbrk(total_size);
  if (block == (void *)-1) {
    pthread_mutex_unlock(&global_malloc_lock);
    return NULL;
  }
  header = block;
  header->s.size = size;
  header->s.is_free = 0;
  header->s.next = NULL;
  if (!head) {
    head = header;
  }
  if (tail) {
    tail->s.next = header;
  }
  tail = header;
  pthread_mutex_unlock(&global_malloc_lock);
  return (void *)(header + 1);
}

// calloc 是分配 num 个长度为 nsize 的连续内存空间
void *calloc(size_t num, size_t nsize) {
  size_t size;
  void *block;
  if (!num || !nsize) {
    return NULL;
  }
  size = num * nsize;
  if (nsize != size / num)
    return NULL;
  block = Malloc(size);
  if (!block)
    return NULL;
  memset(block, 0, size);
  return block;
}

void *realloc(void *block, size_t size) {
  header_t *header;
  void *ret;
  if (!block || !size) {
    return Malloc(size);
  }
  header = (header_t *)block - 1;
  if (header->s.size >= size) {
    return block;
  }
  ret = Malloc(size);
  if (ret) {
    memcpy(ret, block, header->s.size);
    free(block);
  }
  return ret;
}

void print_mem_list() {
  header_t *curr = head;
  printf("head = %p, tail = %p \n", (void *)head, (void *)tail);
  while (curr) {
    printf("addr = %p, size = %zu, is_free=%u, next=%p\n", (void *)curr,
           curr->s.size, curr->s.is_free, (void *)curr->s.next);
    curr = curr->s.next;
  }
}

int main() {
  print_mem_list();

  void *mem100 = Malloc(100);
  print_mem_list();
  void *mem10 = Malloc(10);
  void *memc1010 = calloc(10, 10);

  print_mem_list();

  free(mem10);

  print_mem_list();

  free(memc1010);

  print_mem_list();

  free(mem100);

  print_mem_list();
  return 0;
}
