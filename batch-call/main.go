package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func createServer(rmq *RabbitMQ) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	r.Get("/produce", func(w http.ResponseWriter, r *http.Request) {
		err := rmq.publish100Msg()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		_, err = w.Write([]byte("produced"))
		if err != nil {
			fmt.Println("what?!!")
		}
	})

	err := http.ListenAndServe(":6969", r)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	rmq, err := New()
	if err != nil {
		panic(err)
	}

	go createServer(rmq)

	rmq.Consume()
}
