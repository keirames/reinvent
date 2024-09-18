package main

import (
	"fmt"
	"main/pool"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync/atomic"
	"time"
)

// 1. Accept function to execute
// 2. Function execution has time limit & configurable
// 3. Nums of workers will execute in parallel tasks
// 4. Pending tasks will record in some queue style
// 5. Deadlock or resources exhausted ?
// 6. Graceful shutdown
// 7. Using Redis instead of in-house queue.
// 8. Maintain at least 1 execute by using Redis.
// 9. Performance testing using x ? (x not discovered yet)

func randomNumber(min int, max int) int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	randomNum := r.Intn(max-min+1) + min

	return randomNum
}

func main() {
	f, err := os.Create("mem.prof")
	if err != nil {
		fmt.Println("Could not create memory profile:", err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	p := pool.New()
	var counter int64 = 0

	execTicker := time.NewTicker(100 * time.Microsecond)
	defer execTicker.Stop()

	timeoutTicker := time.After(10 * time.Second)

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
			// Write the memory profile to the file
			if err := pprof.WriteHeapProfile(f); err != nil {
				fmt.Println("Could not write memory profile:", err)
				return
			}
			fmt.Println("Memory profile saved to mem.prof")

			fmt.Printf("Done %v jobs \n", atomic.LoadInt64(&counter))
			return
		}
	}
}
