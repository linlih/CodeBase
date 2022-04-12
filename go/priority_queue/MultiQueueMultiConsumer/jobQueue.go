package main

import (
	"container/list"
	"context"
	"errors"
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
	stopped     bool
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewPriorityQueue(size int) *PriorityQueue {
	p := &PriorityQueue{}
	p.noticeChan = make(chan struct{}, size)
	p.capacity = size
	p.size = 0
	p.priorityIdx = make(map[int]int)
	p.ctx, p.cancel = context.WithCancel(context.Background())
	return p
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
	ctx, _ := context.WithCancel(p.ctx)
	job.SetCtx(ctx)

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
		p.RemoveLeastPriorityJob()
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
	if p.stopped {
		return nil
	}
	for {
		if p.queues[idx].jobList.Len() != 0 {
			break
		}
		idx += 1
	}
	elements := p.queues[idx].jobList.Front()
	job := p.queues[idx].jobList.Remove(elements).(Job)
	//// TODO：如果该优先级的队列为空了，删除是否合理？
	//if p.queues[idx].jobList.Len() == 0 {
	//	p.queues = p.queues[0:idx] //
	//}
	return job
}

func (p *PriorityQueue) RemoveLeastPriorityJob() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.queues) <= 0 || p.queues[0].jobList.Len() <= 0 {
		return errors.New("there is no job in priority queue")
	}
	queue := p.queues[len(p.queues)-1]         // 找到优先级最低的队列
	queue.jobList.Remove(queue.jobList.Back()) // 删除最后一个Job
	return nil
}

func (p *PriorityQueue) CancelAllJob() {
	p.cancel()
}
