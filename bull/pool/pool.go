package pool

import (
	"fmt"
	"sync"
)

var DefaultNumsOfWorkers = 10

type Queue struct {
	items []func()
}

func (q *Queue) Enqueue(f func()) {
	items := make([]func(), len(q.items)+1)
	items[0] = f
	for i, item := range q.items {
		items[i+1] = item
	}
	q.items = items
}

func (q *Queue) GetLength() int {
	return len(q.items)
}

type Pool struct {
	numsOfWorkers     int
	numsOfIdleWorkers int
	queue             Queue
	mu                sync.Mutex
}

func New() *Pool {
	p := new(Pool)
	p.numsOfWorkers = DefaultNumsOfWorkers
	p.numsOfIdleWorkers = DefaultNumsOfWorkers

	return p
}

func (p *Pool) IncNumsOfIdleWorkers() {
	p.mu.Lock()
	p.numsOfIdleWorkers++
	p.mu.Unlock()
}

func (p *Pool) DescNumsOfIdleWorkers() {
	p.mu.Lock()
	p.numsOfIdleWorkers--
	p.mu.Unlock()
}

func (p *Pool) Submit(f func()) error {
	p.queue.Enqueue(f)
	fmt.Println(p.queue)
	fmt.Println(p.queue.GetLength())
	//p.mu.Lock()
	//fmt.Println(p.numsOfIdleWorkers)
	//if p.numsOfIdleWorkers <= 0 {
	//	p.mu.Unlock()
	//	return fmt.Errorf("bad state")
	//}
	//p.mu.Unlock()
	//
	//go func() {
	//	defer func() {
	//		p.IncNumsOfIdleWorkers()
	//	}()
	//
	//	p.DescNumsOfIdleWorkers()
	//
	//	f()
	//	fmt.Println("Job done")
	//}()

	return nil
}
