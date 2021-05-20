package znet

import (
	"fmt"
	"net"
)
import "github.com/vwmin/zinx/ziface"

// 连接模块实现
type Connection struct {
	// 当前连接的TCP套接字
	Conn *net.TCPConn

	// 当前连接ID
	ConnID uint32

	// 当前连接状态
	isClosed bool

	// 通知退出channel
	ExitChan chan bool

	// 该Connection绑定到的Router
	Router ziface.IRouter
}

// 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("ConnID = ", c.ConnID, " Reader is Exit, remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for true {
		// 读数据到缓冲区
		buf := make([]byte, 512)

		if _, err := c.Conn.Read(buf); err != nil {
			fmt.Println("recv buf err, err: ", err)
			continue
		}

		// 得到当前连接的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		go func(request ziface.IRequest) {
			// 从路由中，找到注册绑定的Conn对应的router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 启动当前连接的读业务
	go c.StartReader()

	// todo：启动当前连接的写业务
}

// 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
}

// 获取当前连接绑定的conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接的id
func (c *Connection) GetConnectionID() uint32 {
	return c.ConnID
}

// 获取客户端的ip:port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据给客户端
func (c *Connection) Send(data []byte) error {
	return nil
}

// 连接构造方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) ziface.IConnection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}
