package main

import (
	"fmt"
	"net"
	"time"
)

type Controller struct {
	clientConn net.Conn
	serverConn net.Conn
}

type channelItem struct {
	data []byte
	conn net.Conn
}

func main() {
	clientChannel := make(chan *channelItem)
	serverChannel := make(chan *channelItem)

	controller := Controller{}

	go func() {
		for {
			select {
			case clientItem := <-clientChannel:
				{
					fmt.Println("client:", clientItem.conn.RemoteAddr())
					_, err := controller.serverConn.Write(clientItem.data)
					if err != nil {
						fmt.Println("failed to send to server using controller client")
						panic(err)
					}
				}
			case serverItem := <-serverChannel:
				{
					fmt.Println("server:", serverItem.conn.RemoteAddr())
					fmt.Println("receive item from server", serverItem.data)
					for {
						if controller.clientConn == nil {
							time.Sleep(time.Second * 1)
							fmt.Println("client connection not ready, waiting...")
							continue
						}

						_, err := controller.clientConn.Write(serverItem.data)
						if err != nil {
							fmt.Println("failed to send to server using controller client")
							panic(err)
						}
						fmt.Println("send to server using controller client")
						return
					}
				}
			}
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			fmt.Println("failed to listen:", err)
			panic(err)
		}

		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept:", err)
			panic(err)
		}
		controller.clientConn = conn

		defer conn.Close()

		for {
			data := make([]byte, 1024)
			_, err := conn.Read(data)
			if err != nil {
				fmt.Println("failed to read:", err)
				panic(err)
			}

			clientChannel <- &channelItem{
				data: data,
				conn: conn,
			}
		}
	}()

	go func() {
		conn, err := net.Dial("tcp", "localhost:3307")
		if err != nil {
			fmt.Println("failed to connect to server")
			panic(err)
		}
		controller.serverConn = conn

		defer conn.Close()

		for {
			data := make([]byte, 1024)
			_, err = conn.Read(data)
			if err != nil {
				fmt.Println("failed to read:", err)
				panic(err)
			}

			serverChannel <- &channelItem{
				data: data,
				conn: conn,
			}
		}
	}()

	time.Sleep(time.Minute * 10)
}
