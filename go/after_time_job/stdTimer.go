package main

import (
	"fmt"
	"sync"
	"time"
)

type Timer interface {
	// Stop 阻止定时器启动，停止成功返回 true， 定时器超时或者已经停止返回 false
	// 如果 Stop 返回 false, 那么定时器已经超时，函数 Clock.AfterFunc(d, f) 已经在它自己的 goroutine 开始执行了
	// Stop 是不会等待 f 执行结束再返回的，如果调用者需要知道 f 是否执行结束
	// 那么调用者需要显式地和 f 进行合作
	Stop() bool
	// Reset 修改定时在时间 d 后超时
	// Reset 只能在定时器超时了，或者定时器停止后调用，如果定时器已经知道超时了，那么可以直接调用 Reset
	// 否则需要和函数 Clock.AfterFunc(d, f) 配合
	Reset(d time.Duration)
}

// MonotonicTime 是一个单调的时间读数
type MonotonicTime struct {
	nanoseconds int64
}

// Before 表示当前的时间是否在 u 时间之前
func (mt MonotonicTime) Before(u MonotonicTime) bool {
	return mt.nanoseconds < u.nanoseconds
}

// After 表示当前的时间是否在 u 时间之后
func (mt MonotonicTime) After(u MonotonicTime) bool {
	return mt.nanoseconds > u.nanoseconds
}

// Add 返回增加后的时间读数  mt + d
func (mt MonotonicTime) Add(d time.Duration) MonotonicTime {
	return MonotonicTime{
		nanoseconds: time.Unix(0, mt.nanoseconds).Add(d).Sub(time.Unix(0, 0)).Nanoseconds(),
	}
}

// Sub 返回 mt-u 之间的持续时间，如果结果超过了 Duration 可以存储的最大或者最小值时间，则将返回最大或最小值时间
// 要计算 t-d 的时间，使用 t.Add(-d)
func (mt MonotonicTime) Sub(u MonotonicTime) time.Duration {
	return time.Unix(0, mt.nanoseconds).Sub(time.Unix(0, u.nanoseconds))
}

type Clock interface {
	Now() time.Time
	NowMonotonic() MonotonicTime

	AfterFunc(d time.Duration, f func()) Timer
}

type stdClock struct {
	baseTime        time.Time     `state:"nosave"`
	monotonicOffset MonotonicTime `state:"nosave"`
	monotonicMU     sync.Mutex    `state:"nosave"`
	maxMonotonic    MonotonicTime
}

func NewStdClock() Clock {
	return &stdClock{
		baseTime: time.Now(),
	}
}

var _ Clock = (*stdClock)(nil)

func (*stdClock) Now() time.Time {
	return time.Now()
}

func (s *stdClock) NowMonotonic() MonotonicTime {
	sinceBase := time.Since(s.baseTime)
	if sinceBase < 0 {
		panic(fmt.Sprintf("got negative duration = %s since base time = %s", sinceBase, s.baseTime))
	}

	monotonicValue := s.monotonicOffset.Add(sinceBase)

	s.monotonicMU.Lock()
	defer s.monotonicMU.Unlock()

	if s.maxMonotonic.Before(monotonicValue) {
		s.maxMonotonic = monotonicValue
	}

	return s.maxMonotonic
}

func (*stdClock) AfterFunc(d time.Duration, f func()) Timer {
	return &stdTimer{
		t: time.AfterFunc(d, f),
	}
}

type stdTimer struct {
	t *time.Timer
}

var _ Timer = (*stdTimer)(nil)

func (st *stdTimer) Stop() bool {
	return st.t.Stop()
}

func (st *stdTimer) Reset(d time.Duration) {
	st.t.Reset(d)
}

func NewStdTimer(t *time.Timer) Timer {
	return &stdTimer{
		t: t,
	}
}
