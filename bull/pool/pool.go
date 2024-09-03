package pool

import (
	"errors"
	"fmt"
	"sync"
)

var DefaultNumsOfWorkers = 5
var ErrWorkerIsNil = errors.New("worker is nil")

// Worker the task executor
type Worker struct {
	pool *Pool
	task chan func()
}

func (w *Worker) run() {
	w.pool.numsOfRunningWorkers++
	go func() {
		defer func() {
			w.pool.mu.Lock()
			w.pool.numsOfRunningWorkers--
			w.pool.mu.Unlock()

			w.pool.cond.Signal()
		}()

		for t := range w.task {
			if t == nil {
				return
			}
			t()
			break
		}
	}()
}

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
	capacity             int
	numsOfWorkers        int
	numsOfIdleWorkers    int
	numsOfRunningWorkers int
	queue                Queue
	workers              []Worker
	mu                   sync.Mutex
	cond                 *sync.Cond
}

func New() *Pool {
	p := new(Pool)
	p.numsOfWorkers = DefaultNumsOfWorkers
	p.numsOfIdleWorkers = DefaultNumsOfWorkers
	p.capacity = DefaultNumsOfWorkers
	p.numsOfRunningWorkers = 0
	p.cond = sync.NewCond(&p.mu)

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

func (p *Pool) retrieveWorker() *Worker {
	p.cond.L.Lock()

	if p.capacity > p.numsOfRunningWorkers {
		// still has enough room
		p.cond.L.Unlock()

		// new worker
		w := new(Worker)
		w.pool = p
		w.task = make(chan func())
		w.run()

		return w
	} else {
		// wait for worker released
		// block and wait for available worker
		fmt.Println("block and wait")
		p.cond.Wait()
		fmt.Println("release blocker")
		p.cond.L.Unlock()

		w := new(Worker)
		w.pool = p
		w.task = make(chan func())
		w.run()

		return w
	}

	return nil
}

func (p *Pool) Submit(f func()) error {
	w := p.retrieveWorker()
	if w == nil {
		return ErrWorkerIsNil
	}

	// assign task
	w.task <- f

	return nil
}
