package packet

type Message struct {
	ID uint32
	Payload []byte
}

func NewMessage(id uint32, payload []byte) *Message {
	return &Message{
		ID: id,
		Payload: payload,
	}
}
