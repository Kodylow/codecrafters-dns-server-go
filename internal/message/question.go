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
	name, bytesRead, err := parseDomainName(data, offset)
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
func parseDomainName(data []byte, startOffset int) (string, int, error) {
	if startOffset >= len(data) {
		return "", 0, fmt.Errorf("start offset %d exceeds data length %d", startOffset, len(data))
	}

	var result []string
	var totalBytes int
	var seenPointers = make(map[int]bool)

	currentPos := 0
	for {
		absolutePos := startOffset + currentPos
		if absolutePos >= len(data) {
			return "", 0, fmt.Errorf("incomplete domain name at position %d", absolutePos)
		}

		length := int(data[absolutePos])

		// Handle pointer
		if length&0xC0 == 0xC0 {
			if absolutePos+1 >= len(data) {
				return "", 0, fmt.Errorf("incomplete pointer at position %d", absolutePos)
			}

			pointerOffset := int(uint16(length&0x3F)<<8 | uint16(data[absolutePos+1]))

			if pointerOffset >= len(data) {
				return "", 0, fmt.Errorf("compression pointer offset %d exceeds data length %d", pointerOffset, len(data))
			}

			if seenPointers[pointerOffset] {
				return "", 0, fmt.Errorf("pointer loop detected at offset %d", pointerOffset)
			}
			seenPointers[pointerOffset] = true

			if totalBytes == 0 {
				totalBytes = currentPos + 2 // 2 bytes for compression pointer
			}

			// Get the suffix from the compression pointer
			suffix, _, err := parseDomainName(data, pointerOffset)
			if err != nil {
				return "", 0, fmt.Errorf("failed to parse pointer target at offset %d: %w", pointerOffset, err)
			}

			result = append(result, strings.Split(suffix, ".")...)
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
		end := currentPos + 1 + length
		if startOffset+end > len(data) {
			return "", 0, fmt.Errorf("label at position %d exceeds data bounds: need %d bytes, got %d",
				absolutePos, length, len(data)-absolutePos-1)
		}

		label := string(data[startOffset+currentPos+1 : startOffset+end])
		result = append(result, label)
		currentPos = end
	}

	if totalBytes == 0 {
		totalBytes = currentPos
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
