package server

import (
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/pkg/gotracer"
)

// UDPServer represents a DNS server that listens for DNS queries over UDP.
type UDPServer struct {
	addr string
	conn *net.UDPConn
	log  *gotracer.Logger
	// Add handlers/processors
	messageHandler MessageHandler
}

// New creates a new DNS server instance with default handlers
func New(addr string, log *gotracer.Logger) *UDPServer {
	return &UDPServer{
		addr:           addr,
		log:            log,
		messageHandler: NewDefaultMessageHandler(log),
	}
}

// Start initializes the UDP server and starts listening for DNS queries.
// It resolves the server address, listens for incoming connections, and handles requests.
// Returns an error if the server setup or listening process fails.
func (s *UDPServer) Start() error {
	s.log.Info.Println("Starting UDP server setup...")

	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.conn = conn

	return s.serve()
}

// serve is the main loop for the UDP server.
// It listens for incoming DNS queries, processes them, and sends responses.
// Returns an error if there's an issue with reading from the UDP connection.
func (s *UDPServer) serve() error {
	defer s.conn.Close()
	buf := make([]byte, 512)

	for {
		size, source, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			s.log.Error.Printf("Error receiving data: %v", err)
			return err
		}

		// Handle each request in a goroutine
		go s.handleRequest(buf[:size], source)
	}
}

// handleRequest processes a single DNS request.
// It uses the messageHandler to process the request and send a response.
// Logs any errors that occur during processing.
func (s *UDPServer) handleRequest(data []byte, source *net.UDPAddr) {
	response, err := s.messageHandler.Handle(data)
	if err != nil {
		s.log.Error.Printf("Failed to handle request: %v", err)
		return
	}

	if _, err := s.conn.WriteToUDP(response, source); err != nil {
		s.log.Error.Printf("Failed to send response: %v", err)
	}
}
