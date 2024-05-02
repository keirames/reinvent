package main

import (
	"fmt"
	"net"
	"time"
)

func keepProcessAlive() {
	ticker := time.Tick(time.Hour)

	for {
		<-ticker
	}
}

func createServerWithPeers(addr string, peers []string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Server listen tcp err:", err)
		return
	}
	defer listener.Close()

	for {
		fmt.Println("Waiting incomming connection...")
		// Accept incomming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept incomming connection err:", err)
			continue
		}

		go func() {
			defer conn.Close()

			// Create buffer to read data
			buffer := make([]byte, 1024)

			for {
				fmt.Println("Waiting to read data from client...")
				// Read data from client
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Read data from client err:", err)
					return
				}

				// Process and use the data (here, we'll just print it)
				fmt.Printf("Received: %s\n", buffer[:n])
			}
		}()

		// Connect to peers
		for _, p := range peers {
			conn, err := net.Dial("tcp", p)
			if err != nil {
				fmt.Println("Connect to peer fail", err)
				return
			}
		}
	}
}

func main() {
	addrs := []string{"localhost:5000", "localhost:5001", "localhost:5002"}

	for _, addr := range addrs {
		peers := []string{}
		for _, pAddr := range addrs {
			if pAddr != addr {
				peers = append(peers, pAddr)
			}
		}

		go createServerWithPeers(addr, peers)
	}

	keepProcessAlive()
}
