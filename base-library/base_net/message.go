package base_net

type Message struct {
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

func NewMsgPackage(data []byte) *Message {
	if data != nil {
		return &Message{DataLen: uint32(len(data)), Data: data}
	}

	return &Message{Data: []byte{}}
}

func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
