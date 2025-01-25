package message

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Question represents a DNS question section
type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

// ParseQuestion parses a DNS question from a byte slice starting at the given offset
func ParseQuestion(data []byte, offset int) (Question, int, error) {
	name, bytesRead, err := parseDomainName(data[offset:])
	if err != nil {
		return Question{}, 0, fmt.Errorf("failed to parse domain name at offset %d: %w", offset, err)
	}

	remainingBytes := len(data) - (offset + bytesRead)
	if remainingBytes < 4 {
		return Question{}, 0, fmt.Errorf("insufficient bytes for question type and class: need 4, got %d", remainingBytes)
	}

	qType := binary.BigEndian.Uint16(data[offset+bytesRead : offset+bytesRead+2])
	qClass := binary.BigEndian.Uint16(data[offset+bytesRead+2 : offset+bytesRead+4])

	return Question{
		Name:  name,
		Type:  qType,
		Class: qClass,
	}, bytesRead + 4, nil
}

// parseDomainName parses a DNS domain name from the given byte slice
func parseDomainName(data []byte) (string, int, error) {
	if len(data) == 0 {
		return "", 0, fmt.Errorf("empty data for domain name")
	}

	var result []string
	var totalBytes int
	var seenPointers = make(map[int]bool)

	currentPos := 0
	for {
		if currentPos >= len(data) {
			return "", 0, fmt.Errorf("incomplete domain name at position %d", currentPos)
		}

		length := int(data[currentPos])

		// Handle pointer
		if length&0xC0 == 0xC0 {
			if currentPos+1 >= len(data) {
				return "", 0, fmt.Errorf("incomplete pointer at position %d", currentPos)
			}

			offset := int(uint16(length&0x3F)<<8 | uint16(data[currentPos+1]))

			if seenPointers[offset] {
				return "", 0, fmt.Errorf("pointer loop detected at offset %d", offset)
			}
			seenPointers[offset] = true

			if totalBytes == 0 {
				totalBytes = currentPos + 2
			}

			// Get the suffix from the compression pointer
			suffix, _, err := parseDomainName(data[0:]) // Pass full data buffer
			if err != nil {
				return "", 0, fmt.Errorf("failed to parse pointer target: %w", err)
			}

			parts := strings.Split(suffix, ".")
			result = append(result, parts...)
			break
		}

		// End of domain name
		if length == 0 {
			if totalBytes == 0 {
				totalBytes = currentPos + 1
			}
			break
		}

		// Regular label
		if currentPos+1+length > len(data) {
			return "", 0, fmt.Errorf("incomplete label at position %d: need %d bytes, got %d",
				currentPos, length, len(data)-(currentPos+1))
		}

		label := string(data[currentPos+1 : currentPos+1+length])
		result = append(result, label)
		currentPos += 1 + length
	}

	return strings.Join(result, "."), totalBytes, nil
}

// Encode converts a Question to its wire format
func (q Question) Encode() []byte {
	result := make([]byte, 0)

	// Encode domain name
	parts := strings.Split(q.Name, ".")
	for _, part := range parts {
		result = append(result, byte(len(part)))
		result = append(result, []byte(part)...)
	}
	result = append(result, 0)

	// Type
	typeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBytes, q.Type)
	result = append(result, typeBytes...)

	// Class
	classBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(classBytes, q.Class)
	result = append(result, classBytes...)

	return result
}
