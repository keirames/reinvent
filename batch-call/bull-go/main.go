package main

import (
	"fmt"
	"time"
)

type Pool struct {
	capacity int32

	workers []int
}

func (p *Pool) Submit(task func()) error {
	fmt.Println("received task")

	func() {
		task()
	}()

	return nil
}

func task() {
	for range 10 {
		fmt.Println("task")
	}
}

func main() {
	p := &Pool{
		capacity: 10,
		workers:  []int{},
	}
	for i := range 10 {
		fmt.Println(i)
		err := p.Submit(task)
		if err != nil {
			fmt.Println(err)
		}
	}

	time.Sleep(time.Second * 5)
}
