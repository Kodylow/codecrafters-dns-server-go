package message

import (
	"encoding/binary"
	"fmt"
)

// Question represents a DNS question section
type Question struct {
	Name  []byte
	Type  uint16
	Class uint16
}

// ParseQuestion parses a DNS question from a byte slice starting at the given offset
func ParseQuestion(data []byte, offset int) (Question, int, error) {
	name, bytesRead, err := parseDomainName(data[offset:])
	if err != nil {
		return Question{}, 0, fmt.Errorf("failed to parse domain name: %w", err)
	}

	// Need at least 4 more bytes for Type and Class
	remainingBytes := len(data) - (offset + bytesRead)
	if remainingBytes < 4 {
		return Question{}, 0, fmt.Errorf("insufficient bytes for question type and class")
	}

	// Read Type and Class (2 bytes each)
	qType := binary.BigEndian.Uint16(data[offset+bytesRead : offset+bytesRead+2])
	qClass := binary.BigEndian.Uint16(data[offset+bytesRead+2 : offset+bytesRead+4])

	return Question{
		Name:  name,
		Type:  qType,
		Class: qClass,
	}, bytesRead + 4, nil
}

// parseDomainName parses a DNS domain name from the given byte slice
func parseDomainName(data []byte) ([]byte, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("empty data for domain name")
	}

	var totalBytes int
	var result []byte
	var seenPointers = make(map[int]bool) // Track seen pointers to prevent loops

	currentPos := 0
	for {
		if currentPos >= len(data) {
			return nil, 0, fmt.Errorf("incomplete domain name")
		}

		length := int(data[currentPos])

		// Check if this is a pointer (first two bits are 1)
		if length&0xC0 == 0xC0 {
			if currentPos+1 >= len(data) {
				return nil, 0, fmt.Errorf("incomplete pointer")
			}

			// Calculate pointer offset (remove first two bits and combine with next byte)
			offset := int(uint16(length&0x3F)<<8 | uint16(data[currentPos+1]))

			// Check for pointer loops
			if seenPointers[offset] {
				return nil, 0, fmt.Errorf("pointer loop detected")
			}
			seenPointers[offset] = true

			// If this is the first pointer encountered, save total bytes read
			if totalBytes == 0 {
				totalBytes = currentPos + 2 // 2 bytes for pointer
			}

			// Continue parsing from the offset position
			currentPos = offset
			continue
		}

		if length == 0 {
			result = append(result, 0) // Add null terminator
			if totalBytes == 0 {
				totalBytes = currentPos + 1
			}
			break
		}

		// Check if we have enough bytes for this label
		if currentPos+1+length > len(data) {
			return nil, 0, fmt.Errorf("incomplete label")
		}

		// Append length byte and label content
		result = append(result, data[currentPos:currentPos+1+length]...)
		currentPos += 1 + length
	}

	return result, totalBytes, nil
}

// Encode converts a Question to its wire format
func (q Question) Encode() []byte {
	result := make([]byte, 0, len(q.Name)+4)
	result = append(result, q.Name...)

	typeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBytes, q.Type)
	result = append(result, typeBytes...)

	classBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(classBytes, q.Class)
	result = append(result, classBytes...)

	return result
}
