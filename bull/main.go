package main

import (
	"fmt"
	"main/pool"
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

func main() {
	p := pool.New()

	for i := range 10 {
		err := p.Submit(func() {
			fmt.Printf("I'm number %v executed\n", i)
			//time.Sleep(time.Second * 1)
			fmt.Printf("I'm number %v done\n", i)
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	time.Sleep(time.Second * 5)
}
