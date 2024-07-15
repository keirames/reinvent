package pool

import "sync"

const DEFAULT_NUM = 3

type Pool struct {
	count int
}

func New() *Pool {
	pool := new(Pool)
	pool.count = DEFAULT_NUM

	wg := new(sync.WaitGroup)
	for range pool.count {
		wg.Add(1)
	}

	return pool
}
