package main

import (
	"sync"
	"testing"
)

// Linux 底下使用 perf 工具来测试两个线程切换的开销
// taskset -c 0 perf bench sched pipe -T
// 命令释义：
// taskset: linux 中用于查看或修改一个进程的 CPU 亲和性，也就是将其放在哪个进程之中去
// perf: linux 系统性能分析工具，包含较多的功能，bench 则是内置的 benchmark，sched 是测试调度器的性能
//       pipe是一个子程序，两个线程互发整数1000000，看两个线程的执行时间是否相近，从而测试调度器是否公平
// 以上命令的含义就是利用 taskset 将 perf 程序绑定在同一个 cpu 上执行，同时 perf 执行 pipe 的性能测试得到调度时间
// 得到的结果类似如下：
// # Running 'sched/pipe' benchmark:
// # Executed 1000000 pipe operations between two threads
//
//      Total time: 7.285 [sec]
//
//               7.285557 usecs/op
//                 137257 ops/sec
// 其中 7.28557 us测量的是线程发送和接收消息的时间，所以将这个结果除以 2 可以大致得到线程上下文切换到时间

// 测量上下文切换的开销
// 使用该命令进行测试：go test bench=. -cpu=1，限定只使用一个 cpu
// 得到的结果类似如下：
// BenchmarkContextSwitch   3974917               253.2 ns/op             0 B/op          0 allocs/op
// 切换一个 goroutine 是 ns 级别的，而切换一个线程是 us 级别的
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})
	b.ReportAllocs()
	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin // 启动通知管道
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}
	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin)
	wg.Wait()
}
