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
	result := s.x * s.x
	fmt.Println("SquareJob execute: the result is ", result)
	return nil
}

type AddJob struct {
	*BaseJob
	x int
}

func (a *AddJob) Execute() error {
	result := a.x + a.x
	fmt.Println("AddJob execute: the result is ", result)
	return nil
}
