package pool

import (
	"errors"
	"fmt"
	"sync"
)

var DefaultNumsOfWorkers = 20000
var ErrWorkerIsNil = errors.New("worker is nil")

// Worker the task executor
type Worker struct {
	pool *Pool
	task chan func()
}

func (w *Worker) run() {
	w.pool.mu.Lock()
	w.pool.numsOfRunningWorkers++
	fmt.Println(w.pool.numsOfRunningWorkers)
	w.pool.mu.Unlock()

	go func() {
		defer func() {
			w.pool.cond.Signal()
		}()

		for t := range w.task {
			if t == nil {
				return
			}
			t()

			w.pool.mu.Lock()
			w.pool.numsOfRunningWorkers--
			w.pool.mu.Unlock()

			w.pool.workerCache.Put(w)
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
	workerCache          *sync.Pool
}

func New() *Pool {
	p := new(Pool)
	p.numsOfWorkers = DefaultNumsOfWorkers
	p.numsOfIdleWorkers = DefaultNumsOfWorkers
	p.capacity = DefaultNumsOfWorkers
	p.numsOfRunningWorkers = 0
	p.cond = sync.NewCond(&p.mu)
	p.workerCache = &sync.Pool{}

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

func newWorker(p *Pool) *Worker {
	w, ok := p.workerCache.Get().(*Worker)
	if !ok {
		w := new(Worker)
		w.pool = p
		w.task = make(chan func())
		w.run()

		return w
	}

	p.mu.Lock()
	p.numsOfRunningWorkers++
	p.mu.Unlock()
	return w
}

func (p *Pool) retrieveWorker() (*Worker, error) {
	if p.numsOfRunningWorkers == p.capacity {
		return nil, fmt.Errorf("pool is full")
	}

	w := newWorker(p)
	return w, nil
}

func (p *Pool) Submit(f func()) error {
	w, err := p.retrieveWorker()
	if err != nil {
		return err
	}

	// assign task
	w.task <- f

	return nil
}
