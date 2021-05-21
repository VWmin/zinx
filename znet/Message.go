package znet

type Message struct {

	// 消息id
	id uint32

	// 消息长度
	dataLen uint32

	// 消息内容
	data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{id: id, dataLen: uint32(len(data)), data: data}
}

func NewMessageHead(id uint32, dataLen uint32) *Message {
	return &Message{id: id, dataLen: dataLen}
}

func (m *Message) GetMsgId() uint32 {
	return m.id
}

func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetMsgId(id uint32) {
	m.id = id
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.dataLen = dataLen
}

func (m *Message) SetData(data []byte) {
	m.data = data
}
