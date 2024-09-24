package main

import (
	"fmt"
	"main/pool"
	"math/rand"
	"sync/atomic"
	"time"
)

func randomNumber(min int, max int) int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	randomNum := r.Intn(max-min+1) + min

	return randomNum
}

func main() {
	p := pool.New()
	var counter int64 = 0

	execTicker := time.NewTicker(100 * time.Microsecond)
	defer execTicker.Stop()

	timeoutTicker := time.After(120 * time.Minute)

	go func() {
		for {
			select {
			case <-execTicker.C:
				err := p.Submit(func() {
					time.Sleep(time.Millisecond * time.Duration(randomNumber(1000, 5000)))
					atomic.AddInt64(&counter, 1)
				})
				if err != nil {
					fmt.Println(err)
				}
			case <-timeoutTicker:
				fmt.Printf("Done %v jobs \n", atomic.LoadInt64(&counter))
				return
			}
		}
	}()
}
