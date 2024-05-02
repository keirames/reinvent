package main

import (
	"fmt"
	"time"
)

func keepProcessAlive() {
	ticker := time.Tick(time.Hour)

	for {
		<-ticker
	}
}

func main() {
	addrs := []string{"localhost:5000", "localhost:5001", "localhost:5002"}

	for _, addr := range addrs {
		fmt.Println(addr)
	}

	keepProcessAlive()
}
