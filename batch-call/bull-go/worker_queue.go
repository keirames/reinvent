package main

type WorkerQueue struct {
	data []Worker
}

func (q *WorkerQueue) enqueue(w Worker) {
	q.data = append(q.data, w)
}

func (q *WorkerQueue) dequeue() {
	q.data = q.data[:len(q.data)-1]
}
