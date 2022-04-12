package main

import (
	"errors"
	"fmt"
	"sync"
)

type WorkerManager struct {
	queue       *PriorityQueue
	closeChans  []chan struct{}
	closeChanMu sync.Mutex
	workerNum   int
	stopped     bool
}

func NewWorkerManager(queueSize int, workerNum int) *WorkerManager {
	return &WorkerManager{
		queue:     NewPriorityQueue(queueSize),
		workerNum: workerNum,
		stopped:   false,
	}
}

func (m *WorkerManager) StartWork() error {
	fmt.Println("Start to work")
	for i := 0; i < m.workerNum; i++ {
		m.createWorker()
	}
	return nil
}

func (m *WorkerManager) createWorker() {
	closeChan := make(chan struct{})
	go func(closeChan chan struct{}) {
		for {
			select {
			case <-closeChan:
				fmt.Println("close")
				return
			case <-m.queue.noticeChan:
				job := m.queue.PopJob()
				m.ConsumeJob(job)
			}
		}
	}(closeChan)
	m.closeChanMu.Lock()
	defer m.closeChanMu.Unlock()
	m.closeChans = append(m.closeChans, closeChan)
}

func (m *WorkerManager) CommitJob(job Job) error {
	if !m.stopped {
		m.queue.PushJob(job)
		return nil
	}
	return errors.New("worker manager is stopped")
}

func (m *WorkerManager) ConsumeJob(job Job) {
	defer func() {
		job.Done()
	}()

	job.Execute()
}

// StopWork 关闭worker manager需要关闭所有的worker协程
// 同时要取消所有队列中的已经在跑的任务
func (m *WorkerManager) StopWork() {
	m.queue.CancelAllJob()
	// 关闭所有的worker协程前需要先取消运行中的任务
	for i := 0; i < len(m.closeChans); i++ {
		close(m.closeChans[i])
	}
}
