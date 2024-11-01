package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	port := 3210
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello world"))
		if err != nil {
			panic("why those happen")
		}
	})

	fmt.Println("Server started at http://localhost:" + strconv.Itoa(port))

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		fmt.Printf("Failed to start server %e", err)
		return
	}
}
