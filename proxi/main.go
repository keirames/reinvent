package main

import (
	"fmt"
	"io"
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
	port := 3069

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		url := "http://host.docker.internal:6969" + r.URL.Path
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
		fmt.Println("response", resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Response = resp
		_, err = w.Write(body)
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
