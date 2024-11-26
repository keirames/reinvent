package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func prepareHandshakePacket() ([]byte, error) {
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
	// Null terminator for server version
	packet.WriteByte(0)

	err := binary.Write(packet, binary.LittleEndian, connectionID)
	if err != nil {
		return nil, fmt.Errorf("Failed to write connectionID")
	}

	packet.Write(saltPart1)
	packet.WriteByte(0) // Null terminator for salt

	err = binary.Write(packet, binary.LittleEndian, capabilities)
	if err != nil {
		return nil, fmt.Errorf("Failed to write capabilities")
	}

	packet.WriteByte(charset)

	err = binary.Write(packet, binary.LittleEndian, status)
	if err != nil {
		return nil, fmt.Errorf("Failed to write status")
	}

	// Fill the Capability Flags (Upper 2 Bytes)
	var capabilityFlagsUpper uint16 = 0xFFFF // Example: All bits set
	err = binary.Write(packet, binary.LittleEndian, capabilityFlagsUpper)
	if err != nil {
		return nil, fmt.Errorf("Failed to fill capabilities upper")
	}

	// Fill the Length of Auth-Plugin-Data (1 Byte)
	authPluginDataLength := byte(21) // Example value
	packet.WriteByte(authPluginDataLength)

	// Fill Reserved Bytes (10 Bytes of 0x00)
	reserved := make([]byte, 10)
	packet.Write(reserved)

	// Write filler and extra salt (for simplicity, skipping some details)
	packet.Write(make([]byte, 13))
	packetLength := packet.Len()

	header := make([]byte, 4)
	header[0] = byte(packetLength & 0xFF)         // Least significant byte
	header[1] = byte((packetLength >> 8) & 0xFF)  // Second byte
	header[2] = byte((packetLength >> 16) & 0xFF) // Third byte

	fmt.Println("header", header)
	binary.LittleEndian.PutUint32(header, uint32(packetLength))
	fmt.Println("header after after", header)
	header[3] = 0 // Sequence ID, increment if needed for multiple packets

	return append(header, packet.Bytes()...), nil
}

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
		packet, err := prepareHandshakePacket()
		fmt.Println(packet)
		if err != nil {
			fmt.Println("Failed to prepare package")
			panic(err)
		}

		// Send header and payload
		_, err = conn.Write(packet)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		log.Println("Handling Request")
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Println(buffer)
		// Parse the handshake packet

		pos := 0
		packetLength := buffer[pos]
		fmt.Println("packet length", packetLength)

		// Skip header
		fmt.Println("header", buffer[pos:pos+3])
		pos += 4

		// Skip capabilities
		fmt.Println("client capabilities", buffer[pos:pos+4])
		pos += 4

		// Maximum packet size
		packetSize := buffer[pos : pos+4]
		fmt.Println("packet size", packetSize)
		fmt.Println(binary.LittleEndian.Uint32(packetSize))
		pos += 4

		// Character set
		pos += 1

		// Filler
		pos += 23

		// Username
		oldPos := pos
		fmt.Println("username", buffer[pos:])
		pos += bytes.IndexByte(buffer[pos:], 0x00)
		fmt.Println("username in binary", buffer[oldPos:pos])
		username := string(buffer[oldPos:pos])
		fmt.Println("username", username)
		pos++ // terminator

		serverAddress := "127.0.0.1:3306" // Replace with your server's IP and port

		// Connect to the server
		nconn, err := net.Dial("tcp", serverAddress)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}
		defer conn.Close()
		fmt.Println("Connected to server:", serverAddress)

		// Example: Send a message to the server
		message := []byte{
			176, 0, 0, 1, 141, 162, 26, 0, 0, 0, 0, 0, 45, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 114, 111, 111, 116,
			0, 20, 114, 251, 114, 166, 148, 216, 132, 76, 134, 51, 158, 86, 52, 86,
			99, 214, 51, 121, 231, 128, 109, 121, 115, 113, 108, 0, 109, 121, 115, 113,
			108, 95, 110, 97, 116, 105, 118, 101, 95, 112, 97, 115, 115, 119, 111, 114,
			100, 0, 89, 12, 95, 99, 108, 105, 101, 110, 116, 95, 110, 97, 109, 101, 15,
			71, 111, 45, 77, 121, 83, 81, 76, 45, 68, 114, 105, 118, 101, 114, 3, 95,
			111, 115, 6, 100, 97, 114, 119, 105, 110, 9, 95, 112, 108, 97, 116, 102, 111,
			114, 109, 5, 97, 114, 109, 54, 52, 4, 95, 112, 105, 100, 4, 56, 53, 52, 53,
			12, 95, 115, 101, 114, 118, 101, 114, 95, 104, 111, 115, 116, 9, 49, 50, 55,
			46, 48, 46, 48, 46, 49,
		}
		_, err = nconn.Write(message)
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		mysqlB := make([]byte, 1024)
		_, err = nconn.Read(mysqlB)
		if err != nil {
			fmt.Println("failed to here from mysql")
			panic(err)
		}
		fmt.Println(mysqlB)

		// write to client
		fmt.Println("write to client, pass mysql data")
		_, err = conn.Write(mysqlB)
		if err != nil {
			fmt.Println("Failed to sent to client from mysql")
			panic(err)
		}
	}
}
