package main

import (
	"fmt"
	"goauth/internal/auth"
	"goauth/internal/server"
)

func main() {
	// init gothic
	auth.New()

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
