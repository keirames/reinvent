package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func targetServer() {
	port := 3333
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("headers from target", r.Header)
		_, err := w.Write([]byte("hello world"))
		if err != nil {
			panic("why those happen")
		}
	})

	fmt.Println("Target server started at http://localhost:" + strconv.Itoa(port))

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
	if err != nil {
		fmt.Printf("Failed to start server %e", err)
		return
	}
}

func main() {
	port := 3210

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header
		for name, values := range headers {
			for _, value := range values {
				fmt.Println(name, value)
			}
		}

		client := &http.Client{}
		clientRequest, err := http.NewRequest("GET", "http://localhost:3333", nil)
		if err != nil {
			fmt.Println("Error when build client http", err)
			r.Response.StatusCode = 500
			return
		}
		// clone all headers
		clientRequest.Header = r.Header
		resp, err := client.Do(clientRequest)
		if err != nil {
			fmt.Println("Target rejected", err)
			r.Response.StatusCode = 500
			return
		}
		fmt.Println(resp)

		_, err = w.Write([]byte("hello world"))
		if err != nil {
			panic("why those happen")
		}
	})

	go targetServer()

	fmt.Println("Server started at http://localhost:" + strconv.Itoa(port))

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
	if err != nil {
		fmt.Printf("Failed to start server %e", err)
		return
	}
}
