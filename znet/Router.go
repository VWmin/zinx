package znet

import "github.com/vwmin/zinx/ziface"

/**
  提供IRouter的空实现
  用户可选择性重写Router动作
 */
type BaseRouter struct {}

// 处理业务之前的方法 Hook
func (r *BaseRouter) PreHandle(request ziface.IRequest) {}

// 处理业务之后的方法
func (r *BaseRouter) Handle(request ziface.IRequest) {}

// 处理业务之后的方法
func (r *BaseRouter) PostHandle(request ziface.IRequest) {}