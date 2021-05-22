package ziface

//定义一个服务器接口
type IServer interface {
	//启动
	Start()

	//停止
	Stop()

	//运行
	Server()

	// 给当前的服务注册一个路由，供连接处理使用
	AddRouter(msgId uint32, router IRouter)

	GetConnManager() IConnManager

	SetConnHook(hook IConnHook)
}
