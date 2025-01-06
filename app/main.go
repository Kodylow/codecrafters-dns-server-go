package main

import (
	"fmt"
	"net"
)

// Ensures gofmt doesn't remove the "net" import in stage 1 (feel free to remove this!)
var _ = net.ListenUDP

// main is the entry point for the UDP server application.
// It sets up a UDP listener on a specified address and port,
// then continuously reads incoming data and sends an empty response.
func main() {
	// Resolve the UDP address to listen on.
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	// Listen for incoming UDP packets on the resolved address.
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close() // Ensure the connection is closed when the function exits.

	// Buffer to store incoming data.
	buf := make([]byte, 512)

	for {
		// Read data from the UDP connection.
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		// Convert the received bytes to a string.
		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response to send back to the source.
		response := []byte{}

		// Send the empty response back to the source.
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
