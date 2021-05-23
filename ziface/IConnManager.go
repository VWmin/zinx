package ziface

import "net"

/**
主要目的是控制生成连接的数量
以及server终止时，关闭所有连接
*/
type IConnManager interface {
	Size() int

	Full() bool

	ClearConnections()

	SetConnHook(hook IConnHook)

	AddConnection(conn IConnection)

	DeleteConnection(connID uint32)

	RetrieveConnection(connID uint32) (*IConnection, bool)

	NewConnection(server IServer, conn *net.TCPConn, connID uint32, handler IReqHandler) IConnection
}
