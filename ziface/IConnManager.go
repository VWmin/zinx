package ziface

import "net"

/**
主要目的是控制生成连接的数量
以及server终止时，关闭所有连接
*/
type IConnManager interface {
	NewConnection(server IServer, conn *net.TCPConn, connID uint32, handler IMsgHandler) IConnection

	AddConnection(conn IConnection)

	DeleteConnection(conn IConnection)

	RetrieveConnection(connID uint32) (*IConnection, bool)

	Size() int

	Full() bool

	ClearConnections()

	// fixme: 还是感觉这样实现不太好
	SetConnHook(hook IConnHook)
}
