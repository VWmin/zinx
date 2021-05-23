package znet

import (
	"fmt"
	"github.com/vwmin/zinx/utils"
	"github.com/vwmin/zinx/ziface"
	"sync"
)

/**
消息处理模块的实现
*/
type ReqHandler struct {
	// 存放消息类型对应的router
	apis map[uint32]ziface.IRouter

	// 读写锁
	lock sync.RWMutex

	// 消息（任务队列），BlockingQueue
	taskQueue []chan ziface.IRequest

	// 承载业务的worker数量
	workerSize uint

	// 每个worker拥有的任务队列大小
	workerTaskQueueSize uint
}

func NewReqHandler() *ReqHandler {
	return &ReqHandler{
		apis:                make(map[uint32]ziface.IRouter),
		taskQueue:           make([]chan ziface.IRequest, utils.GlobalObject.WorkerSize),
		workerSize:          utils.GlobalObject.WorkerSize,
		workerTaskQueueSize: utils.GlobalObject.WorkerTaskQueueSize,
	}
}

func (h *ReqHandler) StartRouting() {
	// 开启工作池
	h.StartWorkerPool()
}

func (h *ReqHandler) StopRouting() {
	// todo 如何处理消息队列中未被加载的内容

	// todo 如何分辨go程正在工作或是阻塞

	// todo 如何处理正在处理req的go程

	for i := 0; i < int(h.workerTaskQueueSize); i++ {
		close(h.taskQueue[i])
	}
}

// 根据消息类型 调用Router动作
func (h *ReqHandler) HandleRequest(request ziface.IRequest) {
	h.Submit(request)
}

func (h *ReqHandler) doHandler(request ziface.IRequest) {
	id := request.GetRequestMsg().GetMsgId()
	h.lock.RLock()
	router, ok := h.apis[id]
	h.lock.RUnlock()
	if !ok {
		fmt.Println("msgID: ", id, " have not register to any router.")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

// 为消息类型注册一个Router
func (h *ReqHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 其实如果不允许Server启动后添加Router的话，就没必要上锁
	h.lock.RLock()
	_, ok := h.apis[msgId]
	h.lock.RUnlock()
	if ok {
		fmt.Println("msgID: ", msgId, " have been registered, ignore")
		return
	}

	h.lock.Lock()
	h.apis[msgId] = router
	h.lock.Unlock()
	fmt.Println("router registers to msgID: ", msgId, " success.")
}

// 开启线程池
func (h *ReqHandler) StartWorkerPool() {
	for i := 0; i < int(h.workerSize); i++ {
		h.taskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.WorkerTaskQueueSize)
		// 启动当前worker，阻塞等待消息从channel传递进来
		go h.StartWorker(i, h.taskQueue[i])
	}
}

// 开启工作线程
func (h *ReqHandler) StartWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started...")

	for true {
		select {
		// 消息到达
		case request := <-taskQueue:
			fmt.Println("Worker ID = ", workerID, " dealing request...")
			h.doHandler(request)
		}
	}
}

// 提交任务 轮询模式
func (h *ReqHandler) Submit(request ziface.IRequest) {
	// 提交给消息队列
	// 尝试建立一个conn - worker的关系
	toWorker := request.GetConnection().GetConnectionID() % uint32(h.workerSize)
	h.taskQueue[toWorker] <- request
}
