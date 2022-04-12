package main

import (
	"container/list"
	"sort"
	"sync"
)

type PriorityQueue struct {
	mu          sync.Mutex
	noticeChan  chan struct{}
	queues      []*JobQueue
	priorityIdx map[int]int // 记录的是优先级(key)和queues的下标(value)之间的对应关系
	size        int
	capacity    int
}

type JobQueue struct {
	priority int        // 表示该队列是属于哪个优先级的
	jobList  *list.List // 存储具体的Job
}

func NewJobQueue(priority int) *JobQueue {
	return &JobQueue{
		priority: priority,
		jobList:  list.New(),
	}
}

func (p *PriorityQueue) PushJob(job Job) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var idx int
	var ok bool
	// 判断优先级队列是否已经存在
	if idx, ok = p.priorityIdx[job.Priority()]; !ok {
		idx = p.addPriorityQueue(job.Priority())
	}
	queue := p.queues[idx]
	queue.jobList.PushBack(job)

	p.size++
	if p.size > p.capacity {
		//p.RemoveLeastPriorityJob()
	} else {
		p.noticeChan <- struct{}{}
	}
}

func (p *PriorityQueue) addPriorityQueue(priority int) int {
	n := len(p.queues)
	// 根据优先级进行二分查找，因为队列是通过从高优先级(数值小的)到低优先级(数值大的)进行排列的
	pos := sort.Search(n, func(i int) bool {
		return priority < p.queues[i].priority
	})

	// 更新映射表中优先级和切片索引的对应关系
	for i := pos; i < n; i++ {
		p.priorityIdx[p.queues[i].priority] = i + 1
	}

	tail := make([]*JobQueue, n-pos)
	copy(tail, p.queues[pos:])
	p.queues = append(p.queues[0:pos], NewJobQueue(priority))

	p.queues = append(p.queues, tail...)

	return pos
}

func (p *PriorityQueue) PopJob() Job {
	var idx int = 0
	p.mu.Lock()
	defer p.mu.Unlock()
	for {
		if p.queues[idx].jobList.Len() != 0 {
			break
		}
		idx += 1
	}
	elements := p.queues[idx].jobList.Front()
	return p.queues[idx].jobList.Remove(elements).(Job)
}
