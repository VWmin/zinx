package znet

import (
	"errors"
	"fmt"
	"io"
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
		dataPack := NewDataPack()

		// 读出消息头字节
		msgHeadBuf := make([]byte, dataPack.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), msgHeadBuf); err != nil {
			fmt.Println("recv buf err, ", err)
			continue
		}

		// 拆包为消息对象
		msgHead, err := dataPack.Unpack(msgHeadBuf)
		if err != nil {
			fmt.Println("unpack err, ", err)
			continue
		}

		// 如果有消息体则读出
		var dataBuf []byte
		if msgHead.GetDataLen() > 0 {
			dataBuf = make([]byte, msgHead.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), dataBuf); err != nil {
				fmt.Println("read data err, ", err)
				continue
			}
		}

		// 消息体字节写入消息对象
		msgHead.SetData(dataBuf)

		// 得到当前连接的Request请求数据
		req := Request{
			conn: c,
			msg: msgHead,
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

// 发送消息给客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection already closed while sending msg")
	}
	message := NewMessage(msgId, data)
	packed, err := NewDataPack().Pack(message)
	if err != nil {
		return err
	}
	return c.Send(packed)
}

// 发送数据给客户端
func (c *Connection) Send(data []byte) error {
	_, err := c.Conn.Write(data)
	return err
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
