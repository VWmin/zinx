package ziface

/**
解决粘包问题的封包拆包模块
先读固定长度的head，得到消息类型和长度
再读消息长度，得到具体内容
*/

type IDataPack interface {
	// 获取包的长度的方法
	GetHeadLen() uint32

	// 封包方法
	Pack(msg IMessage) ([]byte, error)

	// 拆包方法
	Unpack(bytes []byte) (IMessage, error)
}
