package znet

import (
	"fmt"
	"github.com/vwmin/zinx/ziface"
)

/**
消息处理模块的实现
*/
type MsgHandler struct {
	// 存放消息类型对应的router
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Apis: make(map[uint32] ziface.IRouter)}
}

// 根据消息类型 调用Router动作
func (h *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	id := request.GetRequestMsg().GetMsgId()
	router, ok := h.Apis[id]
	if !ok {
		fmt.Println("msgID: ", id, " have not register to any router.")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

// 为消息类型注册一个Router
func (h *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := h.Apis[msgId]; ok {
		fmt.Println("msgID: ", msgId, " have been registered, ignore")
		return
	}

	h.Apis[msgId] = router
	fmt.Println("router have been registered to ", msgId)
}
