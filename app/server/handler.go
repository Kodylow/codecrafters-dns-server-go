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
func (h *DefaultMessageHandler) Handle(data []byte) (message.Message, error) {
	header, err := message.ParseHeader(data)
	if err != nil {
		return message.Message{}, fmt.Errorf("failed to parse header: %w", err)
	}

	question, _, err := message.ParseQuestion(data, message.HeaderSize)
	if err != nil {
		return message.Message{}, fmt.Errorf("failed to parse question: %w", err)
	}

	responseHeader := h.buildResponseHeader(header)
	answer := h.buildAnswer(question)

	return message.Message{
		Header:     responseHeader,
		Question:   question,
		Answer:     answer,
		Authority:  []byte{},
		Additional: []byte{},
	}, nil
}

// buildResponseHeader creates a response header based on the request header
func (h *DefaultMessageHandler) buildResponseHeader(header message.Header) message.Header {
	const (
		ResponseBit = 1
		NoError     = 0
		NotImpl     = 4
	)

	response := header
	response.QR = ResponseBit
	response.AA = 0
	response.TC = 0
	response.RA = 0
	response.Z = 0
	response.ANCount = 1

	if header.Opcode == message.StandardQuery {
		response.RCode = NoError
	} else {
		response.RCode = NotImpl
	}

	return response
}

// buildAnswer creates an answer section for the DNS response
func (h *DefaultMessageHandler) buildAnswer(question message.Question) message.Answer {
	const (
		DefaultTTL    = 60
		IPv4Length    = 4
		DefaultIPAddr = 0x08080808 // 8.8.8.8
	)

	return message.Answer{
		Name:   question.Name,
		Type:   question.Type,
		Class:  question.Class,
		TTL:    DefaultTTL,
		Length: IPv4Length,
		RData:  []byte{8, 8, 8, 8},
	}
}
