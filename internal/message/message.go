package message

// Message represents a DNS message
type Message struct {
	Header     Header
	Question   []byte
	Answer     []byte
	Authority  []byte
	Additional []byte
}

// ToBytes converts the Message to a byte slice
func (m *Message) ToBytes() []byte {
	return append(m.Header.ToBytes(), append(m.Question, append(m.Answer, append(m.Authority, m.Additional...)...)...)...)
}
