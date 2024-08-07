package main

import "fmt"

type Pool struct {
	capacity int32

	workers WorkerQueue
}

func (p *Pool) Submit(task func()) error {
	fmt.Println("received task")

	func() {
		task()
	}()

	return nil
}

func NewPool() *Pool {
	p := new(Pool)

	p.workers = WorkerQueue{}

	return p
}
