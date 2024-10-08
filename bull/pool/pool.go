package pool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var DefaultNumOfWorkers = 10
var ErrWorkerIsNil = errors.New("worker is nil")
var ErrPoolFull = errors.New("pool is full")

// Worker the task executor
type Worker struct {
	pool *Pool
	task chan func()
}

func (w *Worker) incRunningWorkers() {
	w.pool.numOfRunningWorkers.Add(1)
}

func (w *Worker) descRunningWorkers() {
	w.pool.numOfRunningWorkers.Add(-1)
}

func (w *Worker) run() {
	w.incRunningWorkers()
	fmt.Println(w.pool.numOfRunningWorkers.Load())

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered in f", r)
			}

			//w.pool.cond.Signal()
			w.descRunningWorkers()
		}()

		for t := range w.task {
			if t == nil {
				return
			}
			t()

			w.descRunningWorkers()

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
	capacity            int
	numOfWorkers        int
	numOfIdleWorkers    int
	numOfRunningWorkers atomic.Int32
	queue               Queue
	workers             []Worker
	mu                  sync.Mutex
	cond                *sync.Cond
	workerCache         *sync.Pool
}

func New() *Pool {
	p := new(Pool)
	p.numOfWorkers = 0
	p.numOfRunningWorkers.Store(0)
	p.capacity = DefaultNumOfWorkers
	p.cond = sync.NewCond(&p.mu)
	p.workerCache = &sync.Pool{}

	return p
}

func (p *Pool) IncNumOfIdleWorkers() {
	p.mu.Lock()
	p.numOfIdleWorkers++
	p.mu.Unlock()
}

func (p *Pool) DescNumOfIdleWorkers() {
	p.mu.Lock()
	p.numOfIdleWorkers--
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

	p.numOfRunningWorkers.Add(1)

	return w
}

func (p *Pool) retrieveWorker() (*Worker, error) {
	if p.numOfRunningWorkers.Load() == int32(p.capacity) {
		return nil, ErrPoolFull
	}
	fmt.Println("get worker success")

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
