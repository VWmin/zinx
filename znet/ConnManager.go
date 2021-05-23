package znet

import (
	"fmt"
	"github.com/vwmin/zinx/utils"
	"github.com/vwmin/zinx/ziface"
	"net"
	"sync"
)

type ConnManager struct {
	conns    map[uint32]ziface.IConnection
	lock     sync.RWMutex
	hook     ziface.IConnHook
	exitChan chan uint32
}

func NewConnManager() *ConnManager {
	m := &ConnManager{
		conns:    make(map[uint32]ziface.IConnection),
		hook:     &ziface.BaseConnHook{},
		exitChan: make(chan uint32),
	}
	go m.listenConnStop()
	return m
}

func (c *ConnManager) AddConnection(conn ziface.IConnection) {
	c.lock.Lock()
	c.conns[conn.GetConnectionID()] = conn
	c.lock.Unlock()

	fmt.Println("Conn ID = ", conn.GetConnectionID(), " add to ConnManager.")
}

func (c *ConnManager) DeleteConnection(connID uint32) {

	c.lock.Lock()
	delete(c.conns, connID)
	c.lock.Unlock()

	fmt.Println("Conn ID = ", connID, " removed from ConnManager.")
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

func (c *ConnManager) listenConnStop() {
	for true {
		select {
		case connID := <-c.exitChan:
			c.DeleteConnection(connID)
		}
	}
}

func (c *ConnManager) NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IReqHandler) ziface.IConnection {
	serverConn := &Connection{
		Conn:             conn,
		ConnID:           connID,
		isClosed:         false,
		Handler:          handler,
		exitChan:         make(chan bool, 1),
		msgChan:          make(chan []byte),
		Server:           server,
		exitChan2Manager: c.exitChan,
		hook:             c.hook,
	}
	serverConn.Start()
	c.AddConnection(serverConn)
	return serverConn
}

func (c *ConnManager) SetConnHook(hook ziface.IConnHook) {
	c.hook = hook
}
