package ziface

type IWorkerPool interface {
	// 开启线程池
	StartWorkerPool()

	// 开启工作线程
	StartWorker(workerID int, taskQueue chan IRequest)

	// 提交任务
	Submit(request IRequest)
}
