package znet

import (
	"errors"
	"fmt"
	"github.com/vwmin/zinx/ziface"
	"net"
)

//IServer的接口实现，定义一个Server的服务模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

// 定义当前客户端连接所绑定的handleAPI，回显功能 fixme：后续需要自定义
func CallBackToClient(conn *net.TCPConn, data []byte, size int) error {
	fmt.Println("[Conn handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:size]); err != nil {
		fmt.Println("write back data err, ", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}

//启动
func (server *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s, Port: %d, is starting \n", server.IP, server.Port)

	// 使用一个Go程承载循环监听业务，避免阻塞在此
	go func() {
		// 1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(server.IPVersion, fmt.Sprintf("%s:%d", server.IP, server.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 2.尝试监听这个地址
		listener, err := net.ListenTCP(server.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", server.IPVersion, " error: ", err)
			return
		}

		fmt.Println("start zinx sever success, [", server.Name, "] listening...")
		var cid uint32
		cid = 0

		// 3.阻塞地等待客户端连接，处理客户端连接业务（读写）
		for true {
			// 如果有client端连接进来，阻塞会返回一个conn句柄
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept error: ", err)
				continue
			}

			// 创建自己的Connection对象
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			//使用conn处理业务
			go dealConn.Start()

		}
	}()

}

//停止
func (server *Server) Stop() {
	//todo：将一些服务器的资源、状态或者创建的链接，停止、回收
}

//运行
func (server *Server) Server() {
	//只进行server的监听功能
	server.Start()

	//todo:做一些启动服务器之后的额外业务

	//阻塞状态，保障上条的异步执行
	select {}
}

//初始化server的方法，返回一个抽象层的Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
