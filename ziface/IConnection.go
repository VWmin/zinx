package ziface

import "net"

// 定义连接的抽象层
type IConnection interface {
	// 启动连接 让当前的连接准备开始工作
	Start()

	// 停止连接 结束当前连接的工作
	Stop()

	// 获取当前连接绑定的conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接的id
	GetConnectionID() uint32

	// 获取客户端的ip:port
	GetRemoteAddr() net.Addr

	// 发送数据给客户端
	Send(data []byte) error
}

// 定义一个处理连接业务的方法 (连接，内容，长度)
type HandleFunc func(*net.TCPConn, []byte, int) error
