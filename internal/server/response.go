package server

import "github.com/codecrafters-io/dns-server-starter-go/internal/message"

// createResponse creates a DNS response header based on the request header.
// It sets the QR (Response/Query) bit to 1, indicating a response.
// The Opcode is copied from the request header, and other fields are set to default values.
// Returns the response header as a byte slice and an error if any.
func createResponse(reqHeader message.Header) ([]byte, error) {
	responseHeader := message.Header{
		ID:      reqHeader.ID,
		QR:      1,
		Opcode:  reqHeader.Opcode,
		AA:      0,
		TC:      0,
		RD:      reqHeader.RD,
		RA:      0,
		Z:       0,
		RCode:   0,
		QDCount: reqHeader.QDCount,
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	return responseHeader.ToBytes(), nil
}
