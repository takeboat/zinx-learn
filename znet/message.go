package znet

type Message struct {
	Id      uint32 // 消息ID
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}
func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetMsgId(id uint32) {
	msg.Id = id
}

func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}
