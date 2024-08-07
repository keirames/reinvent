package main

import (
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
)

type Worker struct {
}

func main() {
	p, err := ants.NewPoolWithFunc(10, func(i interface{}) {
		time.Sleep(time.Millisecond * 500)
		v, ok := i.(int)
		if !ok {
			panic("not ok")
		}

		fmt.Println(v)
	})
	defer p.Release()

	if err != nil {
		fmt.Println("Failed to initiate goroutine pool")
		panic(err)
	}

	// push 1000 jobs
	for i := range 1000 {
		err := p.Invoke(i)
		if err != nil {
			fmt.Println(err)
		}
	}

	time.Sleep(time.Second * 5)
}
