package main

import (
	"fmt"
	"time"
)

func main() {
	workerManager := NewWorkerManager(10, 4)
	go workerManager.StartWork()

	job1 := &SquareJob{
		BaseJob: &BaseJob{
			DoneChan: make(chan struct{}, 1),
			priority: 10,
		},
		x: 5,
	}

	workerManager.CommitJob(job1)

	job2 := &AddJob{
		BaseJob: &BaseJob{
			DoneChan: make(chan struct{}, 1),
			priority: 1,
		},
		x: 10,
	}

	workerManager.CommitJob(job2)

	job3 := &AddJob{
		BaseJob: &BaseJob{
			DoneChan: make(chan struct{}, 1),
			priority: 0,
		},
		x: 100,
	}

	workerManager.CommitJob(job3)

	job4 := &SquareJob{
		BaseJob: &BaseJob{
			DoneChan: make(chan struct{}, 1),
			priority: 4,
		},
		x: 25,
	}

	workerManager.CommitJob(job4)
	//job1.WaitDone()
	//job2.WaitDone()
	//job3.WaitDone()
	//job4.WaitDone()
	time.Sleep(time.Second)
	workerManager.StopWork()
	fmt.Println("Job Done! END")
}
