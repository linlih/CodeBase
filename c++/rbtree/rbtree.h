//
// Created by linlh on 22-6-21.
//

#ifndef C_RBTREE_H
#define C_RBTREE_H

struct rb_node {
    unsigned long rb_parent_color;
    struct rb_node *rb_right;
    struct rb_node *rb_left;
}__attribute__((aligned(sizeof(long))));

struct rb_root {
    struct rb_node *rb_node;
};

#define RB_RED   0
#define RB_BLACK 1

// 利用rb_parent_color存储了父节点的地址
#define rb_parent(r) ((struct rb_node*)((r)->rb_parent_color) & ~3)
#define rb_color(r) ((r)->rb_parent_color & 1)
#define rb_is_red(r) (!rb_color(r))
#define rb_is_black(r) (rb_color(r))
#define rb_set_red(r) do { (r)->rb_parent_color &= ~1; } while(0)
#define rb_set_black(r) do { (r)->rb_parent_color |= 1; } while(0)

static inline void rb_set_parent(struct rb_node *rb, struct rb_node *p) {
    rb->rb_parent_color = (rb->rb_parent_color & 3) | (unsigned long)p;
}

static inline void rb_set_color(struct rb_node* rb, int color) {
    rb->rb_parent_color = (rb->rb_parent_color & ~1) | color;
}

#define RB_ROOT (struct rb_root) {NULL,}

void rb_insert_color(struct rb_node *node, struct rb_root *root);

#endif //C_RBTREE_H
