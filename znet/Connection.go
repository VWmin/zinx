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

	// 通知退出channel，Reader 告知 Writer
	exitChan chan bool

	// 读写通信管道,无缓冲
	msgChan chan []byte

	// 消息分发器
	Handler ziface.IMsgHandler
}

// 连接的读业务
func (c *Connection) StartReader() {
	fmt.Println("ConnID = ", c.ConnID, " Reader Goroutine is running...")

	defer fmt.Println("ConnID = ", c.ConnID, " Reader is Exited, remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for true {
		dataPack := NewDataPack()

		// 读出消息头字节
		msgHeadBuf := make([]byte, dataPack.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), msgHeadBuf); err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client...")
			} else {
				fmt.Println("recv buf err, ", err)
			}
			break
		}

		// 拆包为消息对象
		msgHead, err := dataPack.Unpack(msgHeadBuf)
		if err != nil {
			fmt.Println("unpack err, ", err)
			break
		}

		// 如果有消息体则读出
		var dataBuf []byte
		if msgHead.GetDataLen() > 0 {
			dataBuf = make([]byte, msgHead.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), dataBuf); err != nil {
				fmt.Println("read data err, ", err)
				break
			}
		}

		// 消息体字节写入消息对象
		msgHead.SetData(dataBuf)

		// todo: Go程数量无法控制，改为线程池 （处理业务，占用CPU maybe）
		// 找到对应路由处理方法并执行
		go c.Handler.DoMsgHandler(&Request{
			conn: c,
			msg:  msgHead,
		})
	}
}

// 连接的写业务
func (c *Connection) StartWriter() {
	fmt.Println("ConnID = ", c.ConnID, " Writer Goroutine is running...")

	defer fmt.Println("ConnID = ", c.ConnID, " Writer is Exited, remote addr is ", c.Conn.RemoteAddr().String())

	// 阻塞等待channel消息，写给client
	for true {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Writer send msg err, ", err)
				return
			}
		case <-c.exitChan:
			// 代表Reader已经退出
			return

		}
	}
}

// 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)


	/*（阻塞等待工作，不占用CPU）*/

	// 启动当前连接的读业务
	go c.StartReader()

	// 启动当前连接的写业务
	go c.StartWriter()

}

// 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	// 告知Writer关闭
	c.exitChan <- true

	// 关闭socket连接
	_ = c.Conn.Close()

	// 回收资源
	close(c.exitChan)
	close(c.msgChan)
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

	// 将要发送的消息由管道通知给Writer
	c.msgChan <- packed
	return nil
}

// 连接构造方法
func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Handler:  handler,
		exitChan: make(chan bool, 1),
		msgChan:  make(chan []byte),
	}
}
