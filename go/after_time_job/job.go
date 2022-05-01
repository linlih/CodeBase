package main

import (
	"sync"
	"time"
)

// JobInstance 是一个具体的 Job 实例
// 每次调度 Job 的时候都会创建不同的实例，所以每个 timer 有自己的 earlyReturn 信号
// 这是为了解决当 Job 停止的时候，快速连续重置导致定时器实例的 earlyReturn 信号
// 受到另一个定时器实例影响或被另外一个 Job 实例看到这个信号
//
// 考虑的到以下的场景，当定时器实例共享相同的 earlyReturn 信号
// T1 在有一个锁 L 的情况下 creates stops reset 一个可取消的定时器
// T2，T3，T4 和 T5 都是 goroutine，分别处理 first(A) second(B) third(C) fourth(D) 的定时器实例
// T1: 获得锁 L
// T1: 创建一个新的 Job， 创建实例 A
// T2: 实例 A 启动，获取锁 L 被阻塞
// T1: 尝试停止实例 A，将 earlyReturn 设置为 true
// T1: 调度定时器（创建实例 B）
// T3: 实例 B 启动，获取锁 L 被阻塞
// T1: 尝试停止实例 B，将 earlyReturn 设置为 true
// T1: 调度定时器（创建实例 C）
// T4: 实例 C 启动，获取锁 L 被阻塞
// T1: 尝试停止实例 C，将 earlyReturn 设置为 true
// T1: 调度定时器（创建实例 D）
// T5: 实例 D 启动，获取锁 L 被阻塞
// T1: 释放 L
//
// 现在 T1 释放了 L，这个时候 4 个定时器中任意一个可以获得 L，然后检查 earlyReturn
// 如果 timer 只是检查 earlyReturn 然后什么也不做，实例 D 将永远不会返回，尽管它没有请求停止
// 如果定时器在没有提前返回之前重置了 earlyReturn，那么所有的定时器中只有一个会工作
// 如果 Job 在 resetting 中重置了 earlyReturn，那么所有的定时器都会工作（当有一个期望去做的时候）
// 所以为了解决上诉的问题，每个定时器实例可以用自己的earlyReturn信号
type JobInstance struct {
	timer Timer

	// 用于通知定时器提前返回，当它被停止了之后，同时有其他启动的 timer 想要去获取锁的时候，具体看这个例子：
	// T1: 获得 lock，调用 Cancel()
	// T2: 定时器启动，但是阻塞在获取 lock 中
	// T1: 释放 lock
	// T2: 获得 lock，执行非计划的任务
	//
	// 为了解决这个问题，T1 会检查是否已经有定时器启动了
	// 然后通知这个定时器使用 earlyReturn 来提前返回，这样一旦 T2 获得锁了之后
	// T2 就会查看 earlyReturn 有没有被设置成 true，这样就会做其他的了
	earlyReturn *bool
}

// Stop 会停止 job 实例 j，如果实例 j 还没有启动
// 如果它已经启动了，会阻塞在等待锁，earlyReturn 将会被设置
// 等到它获得锁的时候，就会执行提前返回
func (j *JobInstance) Stop() {
	if j.timer != nil {
		j.timer.Stop()
		*j.earlyReturn = true
	}
}

// Job 代表了计划执行的工作，当这个工作启动的时候，"相关工作" 结束了它也可以被安全关闭
// "相关工作" 定义为有些工作需要被停止，定时器需要在执行工作的时候持有锁
// 注意：拷贝一个 Job 是不安全的，因为 timer 实例创建了一个 Job 实例地址的闭包？
type Job struct {
	//_ sync.NoCopy

	clock Clock

	instance JobInstance

	locker sync.Locker

	fn func()
}

func (j *Job) Cancel() {
	j.instance.Stop()

	j.instance = JobInstance{}
}

// Schedule 在时间 d 后执行调度任务
// 这个函数可以在取消后重新调用，或者是执行完成Job后，重新调用
// Schedule 应该在没有调度，取消，或者执行完成的情况下调用
// 为了保证安全，调用者应该在每次调用 Schedule 之前先调用 Cancel
func (j *Job) Schedule(d time.Duration) {
	earlyReturn := false
	locker := j.locker
	j.instance = JobInstance{
		timer: j.clock.AfterFunc(d, func() {
			// 只有拿到这个锁了，才能够开始执行计划任务
			locker.Lock()
			defer locker.Unlock()
			if earlyReturn {
				// 如果执行到这一步，说明定时器已经启动，当时有其他的 goroutine 在他持有 lock 的时候调用了 Cancel
				// 所以这里直接返回，什么也不做
				earlyReturn = false
				return
			}
			j.fn()
		}),
		earlyReturn: &earlyReturn,
	}
}

// NewJob 返回一个新的 Job，在计划时间到了在它自己的 goroutine 执行
// 在调用 f 之前，它会阻塞住，当 f 返回时候解锁
// var clock StdClock
// var mu sync.Mutex
// message = "foo"
// job := NewJob(&clock, &mu, func() {
//     fmt.Println(message)
// })
// mu.Lock()
// message = "bar"
// mu.Unlock()
// // output: bar
//
// f 不允许尝试去锁住 mu
//
// var clock StdClock
// var mu sync.Mutex
// message = "foo"
// job := NewJob(&clock, &mu, func() {
//     fmt.Println(message)
// })
// mu.Lock()
// job.Cancel()
// mu.Unlock()
func NewJob(c Clock, l sync.Locker, f func()) *Job {
	return &Job{
		clock:  c,
		locker: l,
		fn:     f,
	}
}
