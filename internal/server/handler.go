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
	Handle(data []byte) (message.Message, error)
}

// DefaultMessageHandler is a default implementation of the MessageHandler interface.
// It uses a logger to log information about the DNS message processing.
type DefaultMessageHandler struct {
	log *gotracer.Logger
}

// NewDefaultMessageHandler creates a new instance of DefaultMessageHandler.
// It takes a logger as an argument to enable logging of message handling activities.
func NewDefaultMessageHandler(log *gotracer.Logger) *DefaultMessageHandler {
	return &DefaultMessageHandler{
		log: log,
	}
}

// Handle processes the DNS message contained in the data byte slice.
// It parses the DNS header and logs the header information.
// Returns a byte slice containing the response or an error if parsing fails.
func (h *DefaultMessageHandler) Handle(data []byte) (message.Message, error) {
	header, err := message.ParseHeader(data)
	if err != nil {
		return message.Message{}, fmt.Errorf("failed to parse header: %w", err)
	}

	question, _, err := message.ParseQuestion(data, 12) // Header is 12 bytes
	if err != nil {
		return message.Message{}, fmt.Errorf("failed to parse question: %w", err)
	}

	answer := message.Answer{
		Name:   question.Name,
		Type:   question.Type,
		Class:  question.Class,
		TTL:    60,
		Length: 4,
		RData:  []byte{8, 8, 8, 8},
	}

	// Create response header
	responseHeader := header
	responseHeader.QR = 1 // Set QR bit to 1 for response

	return message.Message{
		Header:     responseHeader,
		Question:   question,
		Answer:     answer,
		Authority:  []byte{},
		Additional: []byte{},
	}, nil
}
