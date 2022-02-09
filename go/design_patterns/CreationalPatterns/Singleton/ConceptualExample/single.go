package main

import (
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

type single struct {
}

var singleInstance *single

// 最开始会有 nil 检查，确保 singleInstance 单例实例在最开始为空。
// 这是为了防止在每次调用的 getInstance 方法时都去执行消耗巨大的锁操作。
// 如果检查不通过，则就意味着 singleIntance 字段已被填充。

// singleInstance 结构体将在锁定期间创建。

// 在获取到锁后还会有一个 nil 检查。这是为了确保即便是有多个协程过了第一次
// 检查，也只能有一个创建单例实例，否则所有协程都会创建自己的单例结构体示例。
func getInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			singleInstance = &single{}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created")
	}

	return singleInstance
}
