package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func randPort() int {
	return rangeIn(1000, 9999)
}

func CreateSimpleServer(ch chan *http.Server, port string) (*http.Server, error) {
	srv := &http.Server{Addr: "localhost:" + port}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	})

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("fail to create a server")
		log.Fatal(err)
		return nil, err
	}

	servers = append(servers, srv)

	return srv, nil
}

var servers []*http.Server

func CreateServer(srv *http.Server) {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	})

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("fail to create a server")
		return
	}

	servers = append(servers, srv)
}

func CreatePingChecker(srv *http.Server) {
	ticker := time.NewTicker(time.Second * 1)

	for {
		<-ticker.C
		fmt.Println("ticker time expired")
		fmt.Println("check heath server", srv.Addr)

		response, err := http.Get(srv.Addr + "/ping")
		if err != nil {
			fmt.Println("server die", err)
			fmt.Println("remove server")

			for idx, s := range servers {
				if s == srv {
					servers = append(servers[:idx], servers[idx+1:]...)
					return
				}
			}

			return
		}

		fmt.Println(response)
	}
}

func CreateServerBackedInPingChecker() {
	srv := &http.Server{Addr: "localhost:" + strconv.Itoa(randPort())}

	go CreateServer(srv)
	go CreatePingChecker(srv)
}

func SpinNewServer() {
	maxCap := 3
	ticker := time.NewTicker(time.Second * 2)

	for {
		<-ticker.C
		if len(servers) < maxCap {
			// create server
		}

	}
}

func CreateAutoScalingGroup() {
	comm := make(chan *http.Server)

	go CreateSimpleServer(comm, "3000")

	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		fmt.Println("ticker time expired")
		fmt.Println("check heath server", servers[0])

		// ping health check
		serverURL := servers[0]
		response, err := http.Get(serverURL.Addr + "/ping")
		if err != nil {
			fmt.Println("server die", err)
			return
		}

		fmt.Println(response)
	}
}
