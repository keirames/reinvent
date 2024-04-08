package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello")
}

func main() {
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://postgres:12345678@localhost:5432/recurrence?sslmode=disable",
	)
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	// ctx := context.Background()
	// req := testcontainers.ContainerRequest{
	// 	Image:        "postgres:latest",
	// 	ExposedPorts: []string{"5432/tcp"},
	// 	WaitingFor:   wait.ForLog("Ready to accept connections"),
	// }

	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	// r.HandleFunc("/create-single-event", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// })
	r.HandleFunc(
		"/create-non-recurrence-event",
		func(w http.ResponseWriter, r *http.Request) {
			var Data struct {
				StartDate string `json:"startDate"`
				EndDate   string `json:"endDate"`
				StartTime string `json:"startTime"`
				EndTime   string `json:"endTime"`
			}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&Data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			startDate, err := time.Parse(time.RFC1123, Data.StartDate)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Println(startDate)

			w.WriteHeader(http.StatusOK)
		},
	).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

	// res, _ := createSingleDayEvent(
	// 	conn,
	// 	time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC),
	// 	time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC),
	// 	time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC),
	// 	time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC),
	// )
	// fmt.Println(*res)
}
