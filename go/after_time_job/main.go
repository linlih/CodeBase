package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	clock := NewStdClock()
	mu := sync.Mutex{}
	hello := func() {
		fmt.Println("hello")
	}
	job := NewJob(clock, &mu, hello)
	job.Schedule(1 * time.Second) // 1 秒后执行该任务

	mu.Lock()
	time.Sleep(2 * time.Second)
	job.Cancel()
	mu.Unlock()
	time.Sleep(time.Second)
	job.Schedule(time.Second)
	time.Sleep(2 * time.Second)
}
