package main

import (
	"fmt"
	"time"
)

type Worker struct {
}

func task() {
	for range 10 {
		fmt.Println("task")
	}
}

func main() {
	p := NewPool()
	fmt.Println(p)

	time.Sleep(time.Second * 5)
}
