package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net"
	"time"
)

type Cache struct {
	ID         int
	clientConn net.Conn
	serverConn net.Conn
}

func main() {
	caches := make([]Cache, 0)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to open tcp server")
		panic(err)
	}

	for {
		fmt.Println("Waiting for tcp connection", caches)
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection", err)
			continue
		}

		cache, err := makeCache(conn)
		if err != nil {
			fmt.Println("Failed to make cache", err)
			continue
		}

		go func() {
			io.Copy(cache.clientConn, cache.serverConn)
		}()

		go func() {
			io.Copy(cache.serverConn, cache.clientConn)
		}()

		caches = append(caches, *cache)
	}
}

func makeCache(clientConn net.Conn) (*Cache, error) {
	// MySQL server address and port
	address := "127.0.0.1:3307"

	// Establish a raw TCP connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Failed to connect to mysql server", err)
		return nil, err
	}

	return &Cache{
		ID:         time.Now().Second(),
		clientConn: clientConn,
		serverConn: conn,
	}, nil
}
