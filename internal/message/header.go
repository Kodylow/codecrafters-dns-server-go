package message

import (
	"encoding/binary"
	"fmt"
)

// Header is the first 12 bytes of a DNS message
// Integers are stored in network byte order (big-endian)
type Header struct {
	// Packet Identifier, 16-bits, A random ID assigned to query packets. Response packets reply with the same ID.
	ID uint16

	// Query/Response Indicator, 1-bit, 0 for question packet, 1 for response packet
	QR uint8

	// Operation Code, 4-bits, Specifies the type of query in the message
	Opcode uint8

	// Authoritative Answer, 1-bit, 0 for no authoritative answer, 1 for authoritative answer
	AA uint8

	// Truncation, 1-bit, 1 if the message is larger than 512 bytes. Always 0 in UDP responses.
	TC uint8

	// Recursion Desired, 1-bit, Sender sets this to 1 if the server should recursively resolve this query, 0 otherwise
	RD uint8

	// Recursion Available, 1-bit, Server sets this to 1 to indicate that recursion is available in the response
	RA uint8

	// Reserved, 3-bits, Used by DNSSEC queries. At inception, this field was reserved for future use.
	Z uint8

	// Response Code, 4-bits, 0 for no error, 1 for format error, 2 for server failure, 3 for name error, 4 for not implemented, 5 for refused, 6-15 for reserved
	RCode uint8

	// Question Count, 16-bits, Number of questions in the Question section
	QDCount uint16

	// Answer Record Count, 16-bits, Number of resource records in the Answer section
	ANCount uint16

	// Authority Record Count, 16-bits, Number of resource records in the Authority section
	NSCount uint16

	// Additional Record Count, 16-bits, Number of resource records in the Additional section
	ARCount uint16
}

// ParseHeader reads a DNS header from a byte slice
func ParseHeader(data []byte) (Header, error) {
	if len(data) < 12 {
		return Header{}, fmt.Errorf("header data too short: got %d bytes, want 12", len(data))
	}

	h := Header{
		ID: binary.BigEndian.Uint16(data[0:2]),
	}

	// Parse flags byte (byte 2)
	flags := data[2]
	h.QR = (flags >> 7) & 0x1
	h.Opcode = (flags >> 3) & 0xF
	h.AA = (flags >> 2) & 0x1
	h.TC = (flags >> 1) & 0x1
	h.RD = flags & 0x1

	// Parse flags byte (byte 3)
	flags = data[3]
	h.RA = (flags >> 7) & 0x1
	h.Z = (flags >> 4) & 0x7
	h.RCode = flags & 0xF

	// Parse counts
	h.QDCount = binary.BigEndian.Uint16(data[4:6])
	h.ANCount = binary.BigEndian.Uint16(data[6:8])
	h.NSCount = binary.BigEndian.Uint16(data[8:10])
	h.ARCount = binary.BigEndian.Uint16(data[10:12])

	return h, nil
}

// ToBytes converts a Header to its wire format representation
func (h *Header) ToBytes() []byte {
	buf := make([]byte, 12)

	// Write ID
	binary.BigEndian.PutUint16(buf[0:2], h.ID)

	// Pack flags (byte 2)
	buf[2] = (h.QR << 7) | (h.Opcode << 3) | (h.AA << 2) | (h.TC << 1) | h.RD

	// Pack flags (byte 3)
	buf[3] = (h.RA << 7) | (h.Z << 4) | h.RCode

	// Write counts
	binary.BigEndian.PutUint16(buf[4:6], h.QDCount)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCount)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCount)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCount)

	return buf
}
