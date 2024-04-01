package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestMain(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	// r.HandleFunc("/create-single-event", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// })

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:6969",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusForbidden, rr.Code)
	}
}
