package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	//TimeOutCtxDemo()
	TestCancelCtxDemo()
}

func TestCancelCtxDemo() {
	ctx, canel := context.WithCancel(context.Background())

	go submission(ctx)
	// do something
	time.Sleep(2 * time.Second)
	//  happened error here ,cancel sub goroutine
	canel()
	fmt.Println("done!")
	for {

	}
}

func submission(ctx context.Context) {
	// 这里也可以do something
	for {
		// 这里也可以do something
		select {
		case <-ctx.Done():
			fmt.Println("being canceled!", ctx.Err())
			return
		default:
			fmt.Println("do something in sub mission")
			time.Sleep(1 * time.Second)
		}
	}
}

func TimeOutCtxDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go handle(ctx, 500*time.Millisecond) // 这个goruntine的执行时间是500ms，所以是可以正常执行完成的
	//go handle(ctx, 1500*time.Millisecond) // 这个goruntine的执行时间是1500ms，这个时候设置的TimeOut执行时间是1000ms，所以这个goroutine收到结束信号无法执行结束
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("Process request with", duration)
	}
}
