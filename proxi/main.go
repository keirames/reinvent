package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	addr := fmt.Sprintf("%s:%s", "localhost", "3333")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Cannot form tcp listener")
		panic(err)
	}

	defer listener.Close()

	for {
		fmt.Println("Start to listen tcp")
		conn, err := listener.Accept()
		fmt.Println("Incomming connection")
		if err != nil {
			fmt.Println("Panic in listener accept phase")
			fmt.Println(err)
			panic(err)
		}

		fmt.Println("prepare data")
		protocolVersion := byte(10)
		serverVersion := "8.3.0"
		connectionID := uint32(12345)
		saltPart1 := []byte("abcd1234")
		capabilities := uint16(0xFFFF)
		charset := byte(33) // utf8_general_ci
		status := uint16(0x0002)

		packet := bytes.NewBuffer(nil)
		packet.WriteByte(protocolVersion)
		packet.WriteString(serverVersion)
		packet.WriteByte(0) // Null terminator for server version
		binary.Write(packet, binary.LittleEndian, connectionID)
		packet.Write(saltPart1)
		packet.WriteByte(0) // Null terminator for salt
		binary.Write(packet, binary.LittleEndian, capabilities)
		packet.WriteByte(charset)
		binary.Write(packet, binary.LittleEndian, status)

		// Write filler and extra salt (for simplicity, skipping some details)
		packet.Write(make([]byte, 13))

		packetLength := len(packet.Bytes())
		fmt.Println("packet", packet.Bytes())

		header := make([]byte, 4)
		header[0] = byte(packetLength & 0xFF)         // Least significant byte
		header[1] = byte((packetLength >> 8) & 0xFF)  // Second byte
		header[2] = byte((packetLength >> 16) & 0xFF) // Third byte

		fmt.Println("header", header)
		binary.LittleEndian.PutUint32(header, uint32(packetLength))
		fmt.Println("header after after", header)
		header[3] = 0 // Sequence ID, increment if needed for multiple packets

		// Send header and payload
		_, err = conn.Write(append(header, packet.Bytes()...))
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		log.Println("Handling Request")
		buffer := make([]byte, 512)
		_, err = conn.Read(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Println(buffer)
		// Parse the handshake packet
		protocolVersion = buffer[0]
		fmt.Println(protocolVersion)
		if protocolVersion != 10 {
			panic("protocol error")
		}

		// Extract server version
		serverVersionEnd := bytes.IndexByte(buffer[1:], 0x00) + 1
		serverVersion = string(buffer[1:serverVersionEnd])
		fmt.Printf("Server Version: %s\n", serverVersion)

		// Extract salt (part 1)
		connectionID = binary.LittleEndian.Uint32(buffer[serverVersionEnd+1:])
		saltPart1 = buffer[serverVersionEnd+5 : serverVersionEnd+13]

		fmt.Printf("Connection ID: %d\n", connectionID)
		fmt.Printf("Salt (part 1): %x\n", saltPart1)

		// Extract remaining salt and server capabilities (skip parsing for brevity)
		// ...
	}
}
