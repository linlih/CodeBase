package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Skipnode 跳表的节点包括了节点值大小Key和实际存放的值Val
type Skipnode struct {
	Key        uint32
	Val        interface{}
	LevelEntry []*Skipnode // 存储每一层指向下一个节点的指针
	Level      int
}

func NewNode(searchKey uint32, value interface{}, createLevel int, maxLevel int) *Skipnode {
	levelEntry := make([]*Skipnode, maxLevel)
	for i := 0; i <= maxLevel-1; i++ {
		levelEntry[i] = nil
	}
	return &Skipnode{Key: searchKey, Val: value, LevelEntry: levelEntry, Level: createLevel}
}

type Skiplist struct {
	Header      *Skipnode
	MaxLevel    int
	Probability float32
	Level       int // 当前跳表的 level
}

const (
	DefaultMaxLevel    int     = 15   // 跳表的最大 level
	DefaultProbability float32 = 0.25 // 默认的概率
)

func NewSkipList() *Skiplist {
	// 头节点的Level是从1开始递增，最大为DefaultMaxLevel=15
	newList := &Skiplist{Header: NewNode(0, "header", 1, DefaultMaxLevel), Level: 1}
	newList.MaxLevel = DefaultMaxLevel
	newList.Probability = DefaultProbability
	return newList
}

func randomP() float32 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Float32()
}

func (s *Skiplist) SetMaxLevel(maxLevel int) {
	s.MaxLevel = maxLevel
}

// RandomLevel 返回的结果是 [0, MaxLevel-1]
func (s *Skiplist) RandomLevel() int {
	level := 1
	for randomP() < s.Probability && level < s.MaxLevel {
		level++
	}
	return level
}

func (s *Skiplist) Search(searchKey uint32) (interface{}, error) {
	currentNode := s.Header
	// 跳表的核心实现，从Level最大的层数开始找，如果当前节点的值小于searchkey，则继续查找下一个节点
	// 否则调到下一层继续查找，最终找到的第0层，也就是完整链表的那一层
	for i := s.Level - 1; i >= 0; i-- {
		for currentNode.LevelEntry[i] != nil && currentNode.LevelEntry[i].Key < searchKey {
			currentNode = currentNode.LevelEntry[i] // 指向下一个节点
		}
	}

	currentNode = currentNode.LevelEntry[0] // 指向完整链表的下一个元素
	if currentNode != nil && currentNode.Key == searchKey {
		return currentNode.Val, nil
	}
	return nil, errors.New("not Found")
}

func (s *Skiplist) Insert(searchKey uint32, value interface{}) {
	updateList := make([]*Skipnode, s.MaxLevel)
	currentNode := s.Header

	// 找到插入的位置
	for i := s.Header.Level - 1; i >= 0; i-- {
		for currentNode.LevelEntry[i] != nil && currentNode.LevelEntry[i].Key < searchKey {
			currentNode = currentNode.LevelEntry[i]
		}
		updateList[i] = currentNode // 记录每一层的最后满足的节点
	}
	currentNode = currentNode.LevelEntry[0]

	// 节点已经存在则更新Val
	if currentNode != nil && currentNode.Key == searchKey {
		currentNode.Val = value
	} else {
		//newLevel := s.RandomLevel() // 生成随机的level高度，实际使用中应该用随机的方式来实现

		// 模拟得到示例的图
		newLevel := 0
		if searchKey == 3 || searchKey == 7 || searchKey == 9 {
			newLevel = 2
		} else if searchKey == 6 {
			newLevel = 3
		} else {
			newLevel = 1
		}
		if newLevel > s.Level { // 如果高度超过了现在的高度
			for i := s.Level + 1; i <= newLevel; i++ { // 超过当前高度的第一个节点都是s.Header
				updateList[i-1] = s.Header
			}
			s.Level = newLevel
			s.Header.Level = newLevel
		}
		newNode := NewNode(searchKey, value, newLevel, s.MaxLevel)
		for i := 0; i <= newLevel-1; i++ {
			newNode.LevelEntry[i] = updateList[i].LevelEntry[i]
			updateList[i].LevelEntry[i] = newNode
		}
	}
}

func (s *Skiplist) Delete(searchKey uint32) error {
	updateList := make([]*Skipnode, s.MaxLevel)
	currentNode := s.Header

	for i := s.Header.Level - 1; i >= 0; i-- {
		for currentNode.LevelEntry[i] != nil && currentNode.LevelEntry[i].Key < searchKey {
			currentNode = currentNode.LevelEntry[i]
		}
		updateList[i] = currentNode
	}

	currentNode = currentNode.LevelEntry[0]
	if currentNode.Key == searchKey {
		for i := 0; i < currentNode.Level; i++ {
			updateList[i].LevelEntry[i] = currentNode.LevelEntry[i]
		}
		for s.Level > 1 && s.Header.LevelEntry[s.Level-1] == nil {
			s.Level--
		}
		currentNode = nil
		return nil
	}
	return errors.New("not found")
}

func (s *Skiplist) DisplayAll() {
	currentNode := s.Header
	// LevelEntry[0] 保存的是所有链表的元素
	for i := s.Level - 1; i >= 0; i-- {
		currentNode = s.Header
		for {
			fmt.Printf("[key:%d][val:%v]->", currentNode.Key, currentNode.Val)
			if currentNode.LevelEntry[i] == nil {
				break
			}
			currentNode = currentNode.LevelEntry[i]
		}
		fmt.Println("")
	}
}
