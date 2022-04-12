package main

import (
	"container/list"
	"fmt"
)

func main() {
	queue := JobQueue{
		jobList:    list.New(),
		noticeChan: make(chan struct{}, 10),
	}
	workerManger := NewWorkerManager(&queue)

	go workerManger.StartWork()

	job := &SquareJob{
		BaseJob: &BaseJob{DoneChan: make(chan struct{}, 1)},
		x:       5,
	}

	queue.PushJob(job)

	job.WaitDone()
	fmt.Println("Job Done, END")
}
