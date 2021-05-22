package ziface

/*
消息管理抽象层
*/
type IMsgHandler interface {
	// 根据消息类型 调用Router动作
	DoMsgHandler(request IRequest)

	// 为消息类型注册一个Router
	AddRouter(msgId uint32, router IRouter)

	// 开启线程池
	StartWorkerPool()

	// 开启工作线程
	StartWorker(workerID int, taskQueue chan IRequest)

	// 提交任务
	Submit(request IRequest)

}
