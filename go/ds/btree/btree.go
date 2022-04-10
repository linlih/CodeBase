package btree

import (
	"sort"
	"sync"
)

// 这个版本的实现中不能保存两个相同的值

const (
	DefaultFreeListSize = 32
)

// Item 表示的是在树中的单一对象
type Item interface {
	// Less 用于比较两个Item之间的大小关系
	// 如果 !a.Less(b) && ！b.Less(a)，表示的就是 a == b，那么在树中只会保存a或者b
	Less(than Item) bool
}

type items []Item
type children []*node

// FreeList 表示的是在btree节点的空闲链表
// 默认每个BTree都有自己的FreeList，但是多个BTree可以共享相同的FreeList
// 两个Btree使用相同的FreeList在并发下写入是安全的
type FreeList struct {
	mu       sync.Mutex // 为了多个Btree并写入的安全
	freelist []*node
}

type copyOnWriteContext struct {
	freelist *FreeList
}

// node 是tree的内部节点
// 它需要保证以下的规则不变：
// * len(children) == 0, len(items) 不受限制
// * len(children) == len(items) + 1
type node struct {
	items    items
	children children
	cow      *copyOnWriteContext
}

var (
	nilItems    = make(items, 16)
	nilChildren = make(children, 16)
)

// NewFreeList 创建一个新的空闲列表
// size表示的free list的空间大小
func NewFreeList(size int) *FreeList {
	return &FreeList{freelist: make([]*node, 0, size)}
}

// 如果free list有空闲的节点，那么就可以直接从list里面取出空闲的返回
func (f *FreeList) newNode() (n *node) {
	f.mu.Lock()
	index := len(f.freelist) - 1
	if index < 0 {
		f.mu.Unlock()
		return new(node)
	}
	n = f.freelist[index]
	f.freelist[index] = nil
	f.freelist = f.freelist[:index]
	f.mu.Unlock()
	return
}

func (f *FreeList) freeNode(n *node) (out bool) {
	f.mu.Lock()
	if len(f.freelist) < cap(f.freelist) {
		f.freelist = append(f.freelist, n)
		out = true
	}
	f.mu.Unlock()
	return
}

type ItemIterator func(i Item) bool

type BTree struct {
	degree int
	length int
	root   *node
	cow    *copyOnWriteContext
}

func New(degree int) *BTree {
	return NewWithFreeList(degree, NewFreeList(DefaultFreeListSize))
}

func NewWithFreeList(degree int, f *FreeList) *BTree {
	if degree <= 1 {
		panic("bad degree")
	}
	return &BTree{
		degree: degree,
		cow:    &copyOnWriteContext{freelist: f},
	}
}

func (s *items) insertAt(index int, item Item) {
	*s = append(*s, nil) // 先插入一个空的内容
	// 这里是不是有问题，如果说index > len(*s)的话，就会导致函数出错
	if index < len(*s) {
		copy((*s)[index+1:], (*s)[index:]) // 将要插入的位置空出来
	}
	(*s)[index] = item
}

func (s *items) removeAt(index int) Item {
	item := (*s)[index]
	copy((*s)[index:], (*s)[index+1:])
	(*s)[len(*s)-1] = nil // 置为空
	*s = (*s)[:len(*s)-1]
	return item
}

// 删除并返回list中最后一个元素
func (s *items) pop() (out Item) {
	index := len(*s) - 1
	out = (*s)[index]
	(*s)[index] = nil
	*s = (*s)[:index]
	return
}

// 从index开始截断
func (s *items) truncate(index int) {
	var toClear items
	*s, toClear = (*s)[:index], (*s)[index:] // 注意切片是左闭右开区间
	for len(toClear) > 0 {
		toClear = toClear[copy(toClear, nilItems):]
	}
}

func (s items) find(item Item) (index int, found bool) {
	i := sort.Search(len(s), func(i int) bool {
		return item.Less(s[i])
	})
	if i > 0 && !s[i-1].Less(item) {
		return i - 1, true
	}
	return i, false
}

func (s *children) insertAt(index int, n *node) {
	*s = append(*s, nil)
	if index < len(*s) {
		copy((*s)[index+1:], (*s)[index:])
	}
	(*s)[index] = n
}

func (s *children) removeAt(index int) *node {
	item := (*s)[index]
	copy((*s)[index:], (*s)[index+1:])
	(*s)[len(*s)-1] = nil // 置为空
	*s = (*s)[:len(*s)-1]
	return item
}

// 删除并返回list中最后一个元素
func (s *children) pop() (out *node) {
	index := len(*s) - 1
	out = (*s)[index]
	(*s)[index] = nil
	*s = (*s)[:index]
	return
}

// 从index开始截断
func (s *children) truncate(index int) {
	var toClear children
	*s, toClear = (*s)[:index], (*s)[index:] // 注意切片是左闭右开区间
	for len(toClear) > 0 {
		toClear = toClear[copy(toClear, nilChildren):]
	}
}

// 没懂这个函数在做什么
func (n *node) mutableFor(cow *copyOnWriteContext) *node {
	if n.cow == cow {
		return n
	}
	out := cow.freelist.newNode()
	if cap(out.items) >= len(n.items) {
		out.items = out.items[:len(n.items)]
	} else {
		out.items = make(items, len(n.items), cap(n.items))
	}
	copy(out.items, n.items)
	if cap(out.children) >= len(n.children) {
		out.children = out.children[:len(n.children)]
	} else {
		out.children = make(children, len(n.children), cap(n.children))
	}
	copy(out.children, n.children)
	return out
}

func (n *node) mutableChild(i int) *node {
	c := n.children[i].mutableFor(n.cow)
	n.children[i] = c
	return c
}

// split 在给定的index位置分割节点，当前节点缩小
// 这个函数会返回在index为i的item
// 一个新的节点包含后面所有的内容
func (n *node) split(i int) (Item, *node) {
	item := n.items[i]
	next := n.cow.freelist.newNode()
	next.items = append(next.items, n.items[i+1:]...)
	n.items.truncate(i)
	if len(n.children) > 0 {
		next.children = append(next.children, n.children[i+1:]...)
		n.children.truncate(i + 1)
	}
	return item, next
}

// maybeSplitChild 检查一个child是否需要split，如果需要split则split
// 返回是否发生了split
func (n *node) maybeSplitChild(i, maxItems int) bool {
	if len(n.children[i].items) < maxItems {
		return false
	}
	first := n.mutableChild(i)
	item, second := first.split(maxItems / 2)
	n.items.insertAt(i, item)
	n.children.insertAt(i+1, second)
	return true
}

func (n *node) insert(item Item, maxItems int) Item {
	i, found := n.items.find(item)
	if found {
		out := n.items[i]
		n.items[i] = item
		return out
	}
	if len(n.children) == 0 {
		n.items.insertAt(i, item)
		return nil
	}
	if n.maybeSplitChild(i, maxItems) {
		inTree := n.items[i]
		switch {
		case item.Less(inTree):
			// 没有改变，需要返回第一个分开的节点
		case inTree.Less(item):
			i++ // 我们要第二个分开的节点
		default:
			out := n.items[i]
			n.items[i] = item
			return out
		}
	}
	return n.mutableChild(i).insert(item, maxItems)
}

func (n *node) get(key Item) Item {
	i, found := n.items.find(key)
	if found {
		return n.items[i]
	} else if len(n.children) > 0 {
		return n.children[i].get(key)
	}
	return nil
}

func min(n *node) Item {
	if n == nil {
		return nil
	}
	for len(n.children) > 0 {
		n = n.children[0]
	}
	if len(n.items) == 0 {
		return nil
	}
	return n.items[0]
}

func max(n *node) Item {
	if n == nil {
		return nil
	}
	for len(n.children) > 0 {
		n = n.children[len(n.children)-1]
	}
	if len(n.items) == 0 {
		return nil
	}
	return n.items[len(n.items)-1]
}

// toRemove 用了表示在 node.remove 调用的时候具体是那个item
type toRemove int

const (
	removeItem = iota
	removeMin
	removeMax
)

func (n *node) remove(item Item, minItems int, typ toRemove) Item {
	var i int
	var found bool
	switch typ {
	case removeMax:
		if len(n.children) == 0 {
			return n.items.pop()
		}
		i = len(n.items)
	case removeMin:
		if len(n.children) == 0 {
			return n.items.removeAt(0)
		}
		i = 0
	case removeItem:
		i, found = n.items.find(item)
		if len(n.children) == 0 {
			if found {
				return n.items.removeAt(i)
			}
			return nil
		}
	default:
		panic("invalid type")
	}
	if len(n.children[i].items) <= minItems {
		return n.grow
	}
}

func (n *node) growChildAndRemove(i int, item Item, minItems int, typ toRemove) Item {
	if i > 0 && len(n.children[i-1].items) > minItems {
		child := n.mutableChild(i)
		stealFrom := n.mutableChild(i - 1)
		stolenItem := stealFrom.items.pop()
		child.items.insertAt(0, n.items[i-1])
		n.items[i-1] = stolenItem

	}
}
