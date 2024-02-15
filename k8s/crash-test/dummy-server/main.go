package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello!"))
		if err != nil {
			fmt.Println("error occur when write", err)
			w.WriteHeader(http.StatusBadGateway)
		} else {
			fmt.Println("Ok")
		}
	})

	router.HandleFunc("/crash", func(w http.ResponseWriter, r *http.Request) {
		log.Fatal("server crashed")
	})

	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server started at localhost:8000")
	log.Fatal(srv.ListenAndServe())
}
