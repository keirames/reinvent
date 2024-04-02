package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "12345678",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start postgres: %s", err)
	}
	defer func() {
		if err := postgresC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop postgres: %s", err)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:6969",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	request, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rr, request)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, rr.Code)
	}
}
