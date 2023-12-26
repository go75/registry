package packet

import (
	"encoding/binary"
	"encoding/json"
)

/**
封包，拆包 模块
*/

func JsonMarshal(id uint32, obj any) ([]byte, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return Pack(NewMessage(id, data)), nil
}

func JsonUnmarshal(data []byte, obj any) (uint32, error) {
	msg := UnPack(data)
	err := json.Unmarshal(msg.Payload, obj)
	if err != nil {
		return 0, err
	}

	return msg.ID, nil
}

// 封包
func Pack(message *Message) []byte {
	// 初始化数据包
	packet := make([]byte, len(message.Payload) + 4)
	
	// 将message的id写入res中
	binary.BigEndian.PutUint32(packet[:4], message.ID)

	// 将message的内容写到packet中
	copy(packet[4:], message.Payload)

	return packet
}

// 拆包
func UnPack(packet []byte) *Message {
	// 解析packet的id
	id := binary.BigEndian.Uint32(packet[:4])
	
	// 组装packet
	msg := &Message{
		ID: id,
		Payload: packet[4:],
	}

	return msg
}