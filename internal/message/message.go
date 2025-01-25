package message

const (
	HeaderSize    = 12
	StandardQuery = 0
)

// Message represents a DNS message
type Message struct {
	Header     Header
	Questions  []Question
	Answers    []Answer
	Authority  []byte
	Additional []byte
}

// Encode converts the Message to a byte slice
func (m *Message) Encode() []byte {
	result := m.Header.Encode()
	for _, q := range m.Questions {
		result = append(result, q.Encode()...)
	}
	for _, a := range m.Answers {
		result = append(result, a.Encode()...)
	}
	return append(result, append(m.Authority, m.Additional...)...)
}
