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
		url := "http://localhost:6969" + r.URL.Path
		method := r.Method

		client := &http.Client{}
		clientRequest, err := http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Println("Error when build client http", err)
			r.Response.StatusCode = 500
			return
		}

		// clone all headers
		clientRequest.Header = r.Header
		clientRequest.Body = r.Body

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

	fmt.Println("Server started at http://localhost:" + strconv.Itoa(port))

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
	if err != nil {
		fmt.Printf("Failed to start server %e", err)
		return
	}
}
