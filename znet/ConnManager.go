package znet

import (
	"fmt"
	"github.com/vwmin/zinx/utils"
	"github.com/vwmin/zinx/ziface"
	"net"
	"sync"
)

type ConnManager struct {
	conns map[uint32]ziface.IConnection
	lock  sync.RWMutex
	hook  ziface.IConnHook
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		conns: make(map[uint32]ziface.IConnection),
		hook:  &BaseConnHook{},
	}
}

func (c *ConnManager) AddConnection(conn ziface.IConnection) {
	c.lock.Lock()
	c.conns[conn.GetConnectionID()] = conn
	c.lock.Unlock()

	fmt.Println("Conn ID = ", conn.GetConnectionID(), " add to ConnManager.")

	c.hook.AfterConnStart(conn)
}

func (c *ConnManager) DeleteConnection(conn ziface.IConnection) {
	c.hook.BeforeConnStop(conn)

	c.lock.Lock()
	delete(c.conns, conn.GetConnectionID())
	c.lock.Unlock()

	fmt.Println("Conn ID = ", conn.GetConnectionID(), " removed from ConnManager.")
}

func (c *ConnManager) RetrieveConnection(connID uint32) (*ziface.IConnection, bool) {
	c.lock.RLock()
	conn, ok := c.conns[connID]
	c.lock.RUnlock()

	return &conn, ok
}

func (c *ConnManager) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.conns)
}

func (c *ConnManager) Full() bool {
	return c.Size() >= utils.GlobalObject.MaxCoon
}

func (c *ConnManager) ClearConnections() {
	size := c.Size()

	c.lock.Lock()
	for connID, conn := range c.conns {
		conn.Stop()
		delete(c.conns, connID)
		fmt.Println("Conn ID = ", connID, " closed while clearing ConnManager.")
	}
	c.lock.Unlock()

	fmt.Println("ConnManager cleared, total = ", size)
}

// fixme: 这算严重耦合吗？
func (c *ConnManager) NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) ziface.IConnection {
	serverConn := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Handler:  handler,
		exitChan: make(chan bool, 1),
		msgChan:  make(chan []byte),
		Server:   server,
	}
	serverConn.Start()
	c.AddConnection(serverConn)
	return serverConn
}

func (c *ConnManager) SetConnHook(hook ziface.IConnHook) {
	c.hook = hook
}
