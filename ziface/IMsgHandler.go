package ziface

/*
消息管理抽象层
*/
type IMessageHandler interface {
	// 根据消息类型 调用Router动作
	DoMsgHandler(request IRequest)

	// 为消息类型注册一个Router
	AddRouter(msgId uint32, router IRouter)
}
