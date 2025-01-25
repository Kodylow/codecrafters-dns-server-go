package message

const (
	HeaderSize    = 12
	StandardQuery = 0
)

// Message represents a DNS message
type Message struct {
	Header     Header
	Question   Question
	Answer     Answer
	Authority  []byte
	Additional []byte
}

// Encode converts the Message to a byte slice
func (m *Message) Encode() []byte {
	return append(m.Header.Encode(), append(m.Question.Encode(), append(m.Answer.Encode(), append(m.Authority, m.Additional...)...)...)...)
}
