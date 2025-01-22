package main

import (
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/pkg/logger"
)

var log *logger.Logger

func main() {
	// Initialize logger
	log = logger.New()

	log.Info.Println("Starting UDP server setup...")

	// Resolve the UDP address to listen on.
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		log.Error.Printf("Failed to resolve UDP address: %v", err)
		return
	}

	// Listen for incoming UDP packets on the resolved address.
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error.Printf("Failed to bind to address: %v", err)
		return
	}
	defer udpConn.Close()

	// Buffer to store incoming data.
	buf := make([]byte, 512)

	for {
		// Read data from the UDP connection.
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Error.Printf("Error receiving data: %v", err)
			break
		}

		// Convert the received bytes to a string.
		receivedData := string(buf[:size])
		log.Info.Printf("Received %d bytes from %s: %s", size, source, receivedData)

		// Create an empty response to send back to the source.
		response := []byte{}

		// Send the empty response back to the source.
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			log.Error.Printf("Failed to send response: %v", err)
		}
	}

	log.Info.Println("UDP server setup complete.")
}
