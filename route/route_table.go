package route

import "registry/packet"

type RouteTable map[uint32] func(data []byte)

func (t RouteTable) Regist(id uint32, fn func(data []byte)) {
	t[id] = fn
}

func (t RouteTable) Process(message *packet.Message) {
	fn := t[message.ID]
	if fn != nil {
		fn(message.Payload)
	}
}