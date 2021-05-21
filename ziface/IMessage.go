package ziface

/**
封装请求数据为Message
*/
type IMessage interface {
	GetMsgId() uint32

	GetDataLen() uint32

	GetData() []byte

	SetMsgId(uint32)

	SetDataLen(uint32)

	SetData([]byte)
}