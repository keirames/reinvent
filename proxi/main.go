package main

import (
	"database/sql"
	"fmt"
	"io"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Controller struct {
	clientConn net.Conn
	serverConn net.Conn
	mysqlConn  net.Conn
}

type channelItem struct {
	data []byte
	conn net.Conn
}

func readBuffer(data []byte) ([]byte, error) {
	// packet length [24 bit]
	pktLen := int(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16)

	if data[pktLen+4] != 0 {
		return data[:pktLen+5], nil
	}

	return data[:pktLen+4], nil
}

func createClient() error {
	db, err := sql.Open("mysql", "root:hipages@tcp(127.0.0.1:8080)/mysql")
	if err != nil {
		fmt.Println("Client: Failed to connect to database", err)
		return err
	}
	defer db.Close()
	fmt.Println("Client: Successfully connected to database")

	// Client will try to ping db every 2 second
	for {
		fmt.Println("Trying to ping database")
		err = db.Ping()
		if err != nil {
			fmt.Println(err)
			fmt.Println("failed to ping db")
		}
		fmt.Println("Successfully pinged database")

		time.Sleep(2 * time.Second)
	}
}

// This server act like a middle man, forward packet from mysql client
// to mysql server
func (c *Controller) createServer() error {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer listener.Close()

	fmt.Println("Listening on port 8080")

	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			fmt.Println("Err occur when accept connection, try again...", err)
			continue
		}

		fmt.Println("Server: Client connected to server")
		c.clientConn = conn
		break
	}

	go func() {
		if c.mysqlConn == nil {
			// Initialize mysql connection
			mysqlConn, err := net.Dial("tcp", "localhost:3307")
			if err != nil {
				fmt.Println("Failed to initialize mysql connection", err)
				panic(err)
			}
			c.mysqlConn = mysqlConn

			// Handshake packet
			p := make([]byte, 1024)
			_, err = mysqlConn.Read(p)
			if err != nil {
				fmt.Println("Failed to read mysql connection", err)
				panic(err)
			}

			fmt.Println("before", p)
			p, err = readBuffer(p)
			if err != nil {
				fmt.Println("Failed to translate buffer", err)
				panic(err)
			}
			fmt.Println("after", p)

			// Send it to client to perform a handshake
			_, err = c.clientConn.Write(p)
			if err != nil {
				fmt.Println("Failed to write handshake to client", err)
				panic(err)
			}
		}
	}()

	for {
		// TODO: find a way more efficient
		p := make([]byte, 1024)
		fmt.Println("Waiting message from client")
		_, err := conn.Read(p)
		if err != nil {
			fmt.Println("Server: Failed to read p from client", err)
			continue
		}
		fmt.Println("Server: Data received", p)

		fmt.Println("before", p)
		p, err = readBuffer(p)
		if err != nil {
			fmt.Println("Failed to translate buffer", err)
			panic(err)
		}
		fmt.Println("after", p)

		_, err = c.mysqlConn.Write(p)
		if err != nil {
			fmt.Println("Server: Failed to write handshake to mysql", err)
			panic(err)
		}
	}
}

func establishClient(conn net.Conn) error {
	fmt.Println("Establishing client")

	mysqlConn, err := net.Dial("tcp", "localhost:3307")
	if err != nil {
		fmt.Println("Failed to establish mysql connection", err)
		return err
	}

	go func() {
		defer mysqlConn.Close()
		defer conn.Close()
		io.Copy(conn, mysqlConn)
	}()

	go func() {
		defer mysqlConn.Close()
		defer conn.Close()
		io.Copy(mysqlConn, conn)
	}()

	return nil
}

func makeServer() error {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Err occur when accept connection, try again...", err)
			continue
		}

		fmt.Println("Server: Client connected to server")
		go establishClient(conn)
	}
}

func main() {
	err := makeServer()
	if err != nil {
		fmt.Println("Failed to make server")
		panic(err)
	}

	// TODO: there is a better way right ?
	time.Sleep(10 * time.Hour)
}

func main1() {
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
						break
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
			fmt.Println("data from client", data)

			validData, err := readBuffer(data)
			if err != nil {
				fmt.Println("failed to read:", err)
				panic(err)
			}

			clientChannel <- &channelItem{
				data: validData,
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
			// data from mysql server
			data := make([]byte, 1024)
			_, err = conn.Read(data)
			if err != nil {
				fmt.Println("failed to read:", err)
				panic(err)
			}
			fmt.Println("data before valid", data)

			validData, err := readBuffer(data)
			if err != nil {
				fmt.Println("failed to read:", err)
				panic(err)
			}
			fmt.Println("send this to client", validData)

			serverChannel <- &channelItem{
				data: validData,
				conn: conn,
			}
		}
	}()

	time.Sleep(time.Minute * 10)
}
