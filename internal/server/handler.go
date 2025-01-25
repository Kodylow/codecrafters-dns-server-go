package server

import (
	"fmt"

	"github.com/codecrafters-io/dns-server-starter-go/internal/message"
	"github.com/codecrafters-io/dns-server-starter-go/pkg/gotracer"
)

// MessageHandler defines an interface for handling DNS messages.
// It provides a method to process incoming data and return a response or an error.
type MessageHandler interface {
	// Handle processes the given byte slice representing a DNS message.
	// It returns a byte slice containing the response or an error if processing fails.
	Handle(data []byte) ([]byte, error)
}

// DefaultMessageHandler is a default implementation of the MessageHandler interface.
// It uses a logger to log information about the DNS message processing.
type DefaultMessageHandler struct {
	log *gotracer.Logger
}

// NewDefaultMessageHandler creates a new instance of DefaultMessageHandler.
// It takes a logger as an argument to enable logging of message handling activities.
func NewDefaultMessageHandler(log *gotracer.Logger) *DefaultMessageHandler {
	return &DefaultMessageHandler{log: log}
}

// Handle processes the DNS message contained in the data byte slice.
// It parses the DNS header and logs the header information.
// Returns a byte slice containing the response or an error if parsing fails.
func (h *DefaultMessageHandler) Handle(data []byte) ([]byte, error) {
	header, err := message.ParseHeader(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	h.log.Info.Printf("Parsed DNS header - ID: %d, QR: %d, Opcode: %d",
		header.ID, header.QR, header.Opcode)

	return createResponse(message.Header(header))
}
