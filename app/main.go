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
		log.Error.Fatal("Failed to resolve UDP address:", err)
	}
	log.Debug.Printf("Resolved UDP address: %s", udpAddr)

	// Rest of your code, replacing the log calls with:
	// log.Info.Printf(...)
	// log.Warn.Printf(...)
	// log.Error.Printf(...)
	// log.Debug.Printf(...)
}
