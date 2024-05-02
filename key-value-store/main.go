package main

import (
	"fmt"
	"main/raft"
	"time"
)

func keepProcessAlive() {
	ticker := time.Tick(time.Hour)

	for {
		<-ticker
	}
}

func main() {
	cm1 := raft.New(1, raft.Leader)
	cm2 := raft.New(2, raft.Follower)
	cm3 := raft.New(3, raft.Follower)

	fmt.Println(cm1, cm2, cm3)

	keepProcessAlive()
}
