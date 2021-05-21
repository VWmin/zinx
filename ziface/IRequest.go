package ziface

/**
 把客户端请求的连接 和 数据 封装在一起
 */
type IRequest interface {

	// 得到当前连接
	GetConnection() IConnection

	// 得到请求数据
	GetRequestMsg() IMessage
}