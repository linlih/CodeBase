package main

import (
	"context"
	"fmt"
)

type Job interface {
	Execute() error
	WaitDone()
	Done()
	Priority() int
	SetCtx(ctx context.Context)
}

type BaseJob struct {
	Err        error
	DoneChan   chan struct{}
	Ctx        context.Context
	cancelFunc context.CancelFunc
	priority   int
}

func (job *BaseJob) Done() {
	close(job.DoneChan)
}

func (job *BaseJob) SetCtx(ctx context.Context) {
	job.Ctx = ctx
}

func (job *BaseJob) WaitDone() {
	select {
	case <-job.DoneChan:
		return
	}
}

func (job *BaseJob) Priority() int {
	return job.priority
}

type SquareJob struct {
	*BaseJob
	x int
}

func (s *SquareJob) Execute() error {
	for {
		select {
		case <-s.Ctx.Done():
			s.Done()
			return nil
		default:
			result := s.x * s.x
			fmt.Println("SquareJob execute: the result is ", result)
			return nil
		}
	}
}

type AddJob struct {
	*BaseJob
	x int
}

func (a *AddJob) Execute() error {
	for {
		select {
		case <-a.Ctx.Done():
			a.Done()
			return nil
		default:
			result := a.x + a.x
			fmt.Println("AddJob execute: the result is ", result)
			return nil
		}
	}
}
