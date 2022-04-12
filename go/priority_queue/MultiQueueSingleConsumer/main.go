package main

import (
	"fmt"
)

func main() {
	var cap int = 10
	queue := &PriorityQueue{
		noticeChan:  make(chan struct{}, cap),
		capacity:    cap,
		priorityIdx: make(map[int]int),
		size:        0,
	}

	workerManager := NewWorkerManager(queue)
	go workerManager.StartWork()

	job1 := &SquareJob{
		BaseJob: &BaseJob{
			DoneChan: make(chan struct{}, 1),
			priority: 10,
		},
		x: 5,
	}

	queue.PushJob(job1)

	job2 := &AddJob{
		BaseJob: &BaseJob{
			DoneChan: make(chan struct{}, 1),
			priority: 1,
		},
		x: 10,
	}

	queue.PushJob(job2)

	job3 := &AddJob{
		BaseJob: &BaseJob{
			DoneChan: make(chan struct{}, 1),
			priority: 0,
		},
		x: 100,
	}

	queue.PushJob(job3)

	job4 := &SquareJob{
		BaseJob: &BaseJob{
			DoneChan: make(chan struct{}, 1),
			priority: 4,
		},
		x: 25,
	}
	queue.PushJob(job4)

	job1.WaitDone()
	job2.WaitDone()
	job3.WaitDone()
	job4.WaitDone()

	fmt.Println("Job Done! END")
}
