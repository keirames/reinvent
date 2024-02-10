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

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server started at 0.0.0.0:8000")
	log.Fatal(srv.ListenAndServe())
}
