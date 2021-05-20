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

	// 当前连接绑定的业务方法
	handleAPI ziface.HandleFunc

	// 通知退出channel
	ExitChan chan bool
}

// 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("ConnID = ", c.ConnID, " Reader is Exit, remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for true {
		// 读数据到缓冲区
		buf := make([]byte, 512)
		size, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err, err: ", err)
			continue
		}

		// 调用当前连接绑定的业务逻辑API

		if err := c.handleAPI(c.Conn, buf, size); err != nil {
			fmt.Println("ConnID = ", c.ConnID, " error while handling, err: ", err)
			break

		}

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
func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI ziface.HandleFunc) ziface.IConnection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		handleAPI: callbackAPI,
		ExitChan:  make(chan bool, 1),
	}
	return c
}
