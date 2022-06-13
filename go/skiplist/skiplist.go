package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Skipnode 跳表的节点包括了节点值大小Key和实际存放的值Val
type Skipnode struct {
	Key     uint32
	Val     interface{}
	Forward []*Skipnode // 存储每一层指向下一个节点的指针
	Level   int
}

func NewNode(searchKey uint32, value interface{}, createLevel int, maxLevel int) *Skipnode {
	forwardEntry := make([]*Skipnode, maxLevel)
	for i := 0; i <= maxLevel-1; i++ {
		forwardEntry[i] = nil
	}
	return &Skipnode{Key: searchKey, Val: value, Forward: forwardEntry, Level: createLevel}
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
	// 否则调到下一层继续查找
	for i := s.Level - 1; i >= 0; i-- {
		for currentNode.Forward[i] != nil && currentNode.Forward[i].Key < searchKey {
			currentNode = currentNode.Forward[i] // 指向下一个节点
		}
	}
	currentNode = currentNode.Forward[0]
	if currentNode != nil && currentNode.Key == searchKey {
		return currentNode.Val, nil
	}
	return nil, errors.New("not Found")
}

func (s *Skiplist) Insert(searchKey uint32, value interface{}) {
	updateList := make([]*Skipnode, s.MaxLevel)
	currentNode := s.Header

	//
	for i := s.Header.Level - 1; i >= 0; i-- {
		for currentNode.Forward[i] != nil && currentNode.Forward[i].Key < searchKey {
			currentNode = currentNode.Forward[i]
		}
		updateList[i] = currentNode
	}
	currentNode = currentNode.Forward[0]

	// 节点已经存在则更新Val
	if currentNode != nil && currentNode.Key == searchKey {
		currentNode.Val = value
	} else {
		newLevel := s.RandomLevel()
		if newLevel > s.Level {
			for i := s.Level + 1; i <= newLevel; i++ {
				updateList[i-1] = s.Header
			}
			s.Level = newLevel
			s.Header.Level = newLevel
		}
		newNode := NewNode(searchKey, value, newLevel, s.MaxLevel)
		for i := 0; i <= newLevel-1; i++ {
			newNode.Forward[i] = updateList[i].Forward[i]
			updateList[i].Forward[i] = newNode
		}
	}
}

func (s *Skiplist) Delete(searchKey uint32) error {
	updateList := make([]*Skipnode, s.MaxLevel)
	currentNode := s.Header

	for i := s.Header.Level - 1; i >= 0; i-- {
		for currentNode.Forward[i] != nil && currentNode.Forward[i].Key < searchKey {
			currentNode = currentNode.Forward[i]
		}
		updateList[i] = currentNode
	}

	currentNode = currentNode.Forward[0]

	if currentNode.Key == searchKey {
		for i := 0; i < currentNode.Level-1; i++ {
			if updateList[i].Forward[i] != nil && updateList[i].Forward[i].Key != currentNode.Key {
				break
			}
			updateList[i].Forward[i] = currentNode.Forward[i]
		}
		for currentNode.Level > 1 && s.Header.Forward[currentNode.Level] == nil {
			currentNode.Level--
		}
		currentNode = nil
		return nil
	}
	return errors.New("not found")
}

func (s *Skiplist) DisplayAll() {
	fmt.Printf("\nhead->")
	currentNode := s.Header
	for {
		fmt.Printf("[key:%d][val:%v]->", currentNode.Key, currentNode.Val)
		if currentNode.Forward[0] == nil {
			break
		}
		currentNode = currentNode.Forward[0]
	}
	fmt.Printf("nil\n")
	fmt.Println("---------------------------------------------------")
	currentNode = s.Header
	for {
		fmt.Printf("[node:%d], val:%v, level:%d ", currentNode.Key, currentNode.Val, currentNode.Level)
		if currentNode.Forward[0] == nil {
			break
		}
		for j := currentNode.Level - 1; j >= 0; j-- {
			fmt.Printf(" fw[%d]:", j)
			if currentNode.Forward[j] != nil {
				fmt.Printf("%d", currentNode.Forward[j].Key)
			} else {
				fmt.Printf("nil")
			}
		}
		fmt.Printf("\n")
		currentNode = currentNode.Forward[0]
	}
	fmt.Printf("\n")
}
