package message

// Message represents a DNS message
type Message struct {
	Header     Header
	Question   []byte
	Answer     []byte
	Authority  []byte
	Additional []byte
}
