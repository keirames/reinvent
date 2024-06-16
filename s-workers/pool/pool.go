package pool

type Pool struct {
	count int
}

func New() *Pool {
	pool := new(Pool)
	pool.count = 3

	return pool
}
