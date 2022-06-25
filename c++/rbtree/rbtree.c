//
// Created by linlh on 22-6-21.
//

#include "rbtree.h"

static void __rb_rotate_left(struct rb_node *node, struct rb_root *root) {
    struct rb_node *right = node->rb_right;
    struct rb_node *parent = rb_parent(node);

    /* 赋值表达式返回的是左边的变量，下面这句话等价于：
     * node->rb_right = right->rb_left;
     * if (node->rb_right != NULL)
     * */
    if ((node->rb_right = right->rb_left))
        rb_set_parent(right->rb_left, node);
    right->rb_left = node;

    rb_set_parent(right, parent);

    if (parent) {
        if (node == parent->rb_left)
            parent->rb_left = right;
        else
            parent->rb_right = right;
    } else
        root->rb_node = right;
    rb_set_parent(node, right);
}

static void __rb_rotate_right(struct rb_node *node, struct rb_root *root) {
    struct rb_node *left = node->rb_left;
    struct rb_node *parent = rb_parent(node);

    if ((node->rb_left = left->rb_right))
        rb_set_parent(left->rb_right, node);
    left->rb_right = node;

    rb_set_parent(left, parent);

    if (parent) {
        if (node == parent->rb_right)
            parent->rb_right = left;
        else
            parent->rb_left = left;
    } else
        root->rb_node = left;
    rb_set_parent(node, left);
}

// 根据红黑树的插入逻辑，进行不同的旋转和染色
void rb_insert_color(struct rb_node *node, struct rb_root *root) {
    struct rb_node *parent, *gparent;

    while((parent = rb_parent(node) && rb_is_red(parent))) {
        gparent = rb_parent(parent);
        if (parent == gparent->rb_left) {
            {
                struct rb_node *uncle = gparent->rb_right;
                // 当uncle节点是红色和parent都是红色的时候
                if (uncle && rb_is_red(uncle)) {
                    rb_set_black(uncle);
                    rb_set_black(parent);
                    node = gparent;
                    continue;
                }
            }
            // uncle节点是黑色
            if (parent->rb_right == node) {
                register struct rb_node *tmp;
                __rb_rotate_left(parent, root);
                tmp = parent;
                parent = node;
                node = tmp;
            }
            rb_set_black(parent);
            rb_set_red(gparent);
            __rb_rotate_right(gparent, root);
        } else {
            {
                struct rb_node *uncle = gparent->rb_left;
                // 当uncle节点是红色和parent都是红色的时候
                if (uncle && rb_is_red(uncle)) {
                    rb_set_black(uncle);
                    rb_set_black(parent);
                    node = gparent;
                    continue;
                }
            }

            if (parent->rb_left == node) {
                register struct rb_node *tmp;
                __rb_rotate_right(parent, root);
                tmp = parent;
                parent = node;
                node = tmp;
            }

            rb_set_black(parent);
            rb_set_red(gparent);
            __rb_rotate_left(gparent, root);
        }
    }
    rb_set_black(root->rb_node);
}

static void __rb_erase_color(struct rb_node *node, struct rb_node *parent, struct rb_root *root) {
    struct rb_node *other;
    while ((!node || rb_is_black(node)) && node != root->rb_node) {
        // 如果节点node是父亲parent的左节点
        if (parent->rb_left == node) {
            // other节点是node的兄弟节点sibling
            other = parent->rb_right;
            if (rb_is_red(other)) { // sibling节点是红色的
                rb_set_black(other);
                rb_set_red(parent);
                __rb_rotate_left(parent, root); // 执行左旋操作
                other = parent->rb_right;
            }
            if ((!other->rb_left || rb_is_black(other->rb_left)) &&
                    (!other->rb_right || rb_is_black(other->rb_right)))
            ) {
                rb_set_red(other);
                node = parent;
                parent = rb_parent(node);
            } else {
                if (!other->rb_right || rb_is_black(other->rb_right)) {
                    rb_set_black(other->rb_left);
                    rb_set_red(other);
                    __rb_rotate_right(other, root);
                    other = parent->rb_right;
                }
                rb_set_color(other, rb_color(parent));
                rb_set_black(parent);
                rb_set_black(other->rb_right);
                __rb_rotate_left(parent, root);
                node = root->rb_node;
                break;
            }
        } else {
            other = parent->rb_left;
            if (rb_is_red(other)) {
                rb_set_black(other);
                rb_set_red(parent);
                __rb_rotate_right(parent, root);
                other = parent->rb_left;
            }
            if ((!other->rb_left || rb_is_black(other->rb_left)) &&
                    (!other->rb_right || rb_is_black(other->rb_right))) {
                rb_set_red(other);
                node = parent;
                parent = rb_parent(node);
            } else {
                if (!other->rb_left || rb_is_black(other->rb_left)) {
                    rb_set_black(other->rb_right);
                    rb_set_red(other);
                    __rb_rotate_left(other, root);
                    other = parent->rb_left;
                }
                rb_set_color(other, rb_color(parent));
                rb_set_black(parent);
                rb_set_black(other->rb_left);
                __rb_rotate_right(parent, root);
                node = root->rb_node;
                break;
            }
        }
        if (node)
            rb_set_black(node);
    }
}

void rb_erase(struct rb_node *node, struct rb_root *root) {

}