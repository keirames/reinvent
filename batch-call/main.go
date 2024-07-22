package main

func main() {
	rmq, err := New()
	if err != nil {
		panic(err)
	}

	rmq.Consume()
}
