package main

import (
	"fmt"
	"runtime"
	"sync"
)

// 这个函数用来测量一个 goroutine 需要占用的内存空间大小是多少
// 使用的方法是启动大量的 goroutine，这些 goroutine 都是空的，并且阻塞等待
// 创建了一个缓冲区为 0 的 chan，每个 goroutine 阻塞读取这个 chan
func main() {
	// 定义了读取内存数值的函数
	memConsumed := func() uint64 {
		runtime.GC() // GC 后计算内存大小
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }

	const numGoroutines = 1e5 // 创建大量的 goroutine，使用大数定律接近一个 goroutine 的大小
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()

	// 得到的结果是：2.607kb，不同的电脑可能执行结果稍有差异
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}
