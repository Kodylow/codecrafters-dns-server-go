package message

import (
	"encoding/binary"
	"strings"
)

// Answer represents a DNS answer
type Answer struct {
	// Name is the domain name of the answer, encoded as a sequence of labels
	Name []byte
	// Type, 2 bytes, 0x0001 for A record, 0x0005 for CNAME, etc.
	Type uint16
	// Class, 2 bytes, usually set to 0x0001 for IN, 0x0002 for CH, etc.
	Class uint16
	// TTL, 4 bytes, the duration in seconds a record can be cached before requerying
	TTL uint32
	// Length, 2 bytes, length of the RDATA field
	Length uint16
	// RDATA, variable length, data specific to the record type
	RData []byte
}

func NewAnswer(domain string) *Answer {
	return &Answer{
		Name:   encodeDomainName(domain),
		Type:   1, // A record
		Class:  1, // IN class
		TTL:    60,
		Length: 4, // IPv4 address length
		RData:  []byte{8, 8, 8, 8},
	}
}

func (a *Answer) Encode() []byte {
	result := make([]byte, 0)

	// Name
	result = append(result, a.Name...)

	// Type (2 bytes)
	typeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBytes, a.Type)
	result = append(result, typeBytes...)

	// Class (2 bytes)
	classBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(classBytes, a.Class)
	result = append(result, classBytes...)

	// TTL (4 bytes)
	ttlBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ttlBytes, a.TTL)
	result = append(result, ttlBytes...)

	// Length (2 bytes)
	lengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBytes, a.Length)
	result = append(result, lengthBytes...)

	// RDATA
	result = append(result, a.RData...)

	return result
}

func encodeDomainName(domain string) []byte {
	var result []byte
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		result = append(result, byte(len(part)))
		result = append(result, []byte(part)...)
	}
	result = append(result, 0) // null terminator
	return result
}
