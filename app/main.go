package main

import (
	"log"
	"net"
)

// main is the entry point for the UDP server application.
// It sets up a UDP listener on a specified address and port,
// then continuously reads incoming data and sends an empty response.
func main() {
	// Set up logger with timestamp, file:line, and severity level
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)

	// Optionally log to a file instead of stdout
	// f, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	//     log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	// log.SetOutput(f)

	log.Println("Starting UDP server setup...")

	// Resolve the UDP address to listen on.
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		log.Fatal("Failed to resolve UDP address:", err)
	}
	log.Printf("Resolved UDP address: %s", udpAddr)

	// Listen for incoming UDP packets on the resolved address.
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Failed to bind to address:", err)
	}
	defer udpConn.Close() // Ensure the connection is closed when the function exits.

	log.Printf("Server listening on %s", udpAddr)

	// Buffer to store incoming data.
	buf := make([]byte, 512)

	for {
		log.Println("Waiting for incoming data...")

		// Read data from the UDP connection.
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error receiving data: %v", err)
			continue
		}

		// Convert the received bytes to a string.
		receivedData := string(buf[:size])
		log.Printf("Received %d bytes from %s: %s", size, source, receivedData)

		// Create an empty response to send back to the source.
		response := []byte{}

		// Send the empty response back to the source.
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			log.Printf("Failed to send response: %v", err)
		} else {
			log.Printf("Sent response to %s", source)
		}
	}
}
