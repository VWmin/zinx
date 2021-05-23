package ziface

/**
  路由抽象接口
  路由内的数据都是IRequest
  将指令映射到具体的方法
 */
type IRouter interface {

	// 处理业务之前的方法 Hook
	PreHandle(request IRequest)

	// 处理业务之后的方法
	Handle(request IRequest)

	// 处理业务之后的方法
	PostHandle(request IRequest)
}


/**
  提供IRouter的空实现
  用户可选择性重写Router动作
*/
type BaseRouter struct {}

// 处理业务之前的方法 Hook
func (r *BaseRouter) PreHandle(request IRequest) {}

// 处理业务之后的方法
func (r *BaseRouter) Handle(request IRequest) {}

// 处理业务之后的方法
func (r *BaseRouter) PostHandle(request IRequest) {}
