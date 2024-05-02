package raft

import (
	"fmt"
	"sync"
	"time"

	"math/rand/v2"
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
	Dead
)

func (s State) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	case Dead:
		return "Dead"
	default:
		panic("unreachable")
	}
}

// ConsensusModule (CM) implements a single node of Raft consensus.
type ConsensusModule struct {
	// mu protects concurrent access to a CM.
	mu sync.Mutex

	// id is the server ID of this CM.
	id int

	// Persistent Raft state on all servers
	currentTerm int

	// Volatile Raft state on all servers
	state State
}

func New(id int, state State) *ConsensusModule {
	cm := new(ConsensusModule)
	cm.id = id
	cm.state = state

	return cm
}

// generates a pseudo-random election timeout duration.
func electionTimeout() time.Duration {
	return time.Duration(150+rand.IntN(150)) * time.Millisecond
}

func runElectionTimer() {}

// Also become an candidate
func (cm *ConsensusModule) startElection() {
	cm.state = Candidate
	cm.currentTerm += 1
	savedCurrentTerm := cm.currentTerm
	fmt.Printf("becomes Candidate (currentTerm=%v)", savedCurrentTerm)

	// votesReceived := 1
	// Send RequestVote RPCs to all other servers concurrently.
}
