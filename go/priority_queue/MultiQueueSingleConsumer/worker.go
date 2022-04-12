package main

import "fmt"

type WorkerManager struct {
	queue     *PriorityQueue
	closeChan chan struct{}
}

func NewWorkerManager(queue *PriorityQueue) *WorkerManager {
	return &WorkerManager{
		queue:     queue,
		closeChan: make(chan struct{}, 1),
	}
}

func (m *WorkerManager) StartWork() error {
	fmt.Println("Start to work")
	for {
		select {
		case <-m.closeChan:
			fmt.Println("closing")
			return nil
		case <-m.queue.noticeChan:
			job := m.queue.PopJob()
			m.ConsumeJob(job)
		}
	}
}

func (m *WorkerManager) ConsumeJob(job Job) {
	defer func() {
		job.Done()
	}()

	job.Execute()
}

func (m *WorkerManager) StopWork() {
	close(m.closeChan)
}
