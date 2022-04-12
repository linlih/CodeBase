package main

import (
	"container/list"
	"sync"
)

type JobQueue struct {
	mu         sync.Mutex    // 需要并发安全
	jobList    *list.List    // list是golang实现的双向队列，其中每个元素保存的是一个job
	noticeChan chan struct{} // 入队一个job就往chan放入一个消息，供消费者消费
}

// PushJob 队列的Push操作
// 这里有个问题：就是说为什么不把Job直接放到Channel中，让消费者进行消费就好了
// 这是因为对于单队列是可以的，但是这里是为了实现优先级的多队列，所以即使是单队列，也是放在list中
func (q *JobQueue) PushJob(job Job) {
	q.jobList.PushBack(job)
	q.noticeChan <- struct{}{}
}

func (q *JobQueue) PopJob() Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.jobList.Len() == 0 {
		return nil
	}
	elements := q.jobList.Front()
	return q.jobList.Remove(elements).(Job)
}

func (q *JobQueue) WaitJob() <-chan struct{} {
	return q.noticeChan
}
