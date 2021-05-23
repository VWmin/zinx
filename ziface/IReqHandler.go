package ziface

/*
消息管理抽象层
*/
type IReqHandler interface {
	// 开始路由，开辟工作资源
	StartRouting()

	// 停止路由，清理工作资源
	StopRouting()

	// 根据消息类型 调用Router动作
	HandleRequest(request IRequest)

	// 为消息类型注册一个Router
	AddRouter(msgId uint32, router IRouter)
}
