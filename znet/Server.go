package znet

import (
	"fmt"
	"github.com/vwmin/zinx/utils"
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

	// 多路由管理器
	ReqToRouter ziface.IReqHandler

	// 连接管理器
	ConnManager ziface.IConnManager
}

//启动
func (server *Server) Start() {

	fmt.Printf("[Start] Server %s Listenner at IP: %s, Port: %d, is starting \n", server.Name, server.IP, server.Port)

	// 启动req处理线程
	server.ReqToRouter.StartRouting()

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

			// 判断是否超出ConnManager容量
			if server.ConnManager.Full() {
				_ = conn.Close()
				fmt.Println("a conn was closed due to conn size overflow.")
				// todo: 响应给客户端
				continue
			}

			// 创建自己的Connection对象
			server.ConnManager.NewConnection(server, conn, cid, server.ReqToRouter)
			cid++

		}
	}()

}

//停止
func (server *Server) Stop() {
	// 将一些服务器的资源、状态或者创建的链接，停止、回收
	fmt.Println("[Stop] Server ", server.Name)
	server.ConnManager.ClearConnections()
	server.ReqToRouter.StopRouting()
}

//运行
func (server *Server) Server() {
	//只进行server的监听功能
	server.Start()

	//todo:做一些启动服务器之后的额外业务

	//阻塞状态，保障上条的异步执行
	// 主线程结束 开启的go程也会结束
	select {}
}

// 添加一个路由
func (server *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	server.ReqToRouter.AddRouter(msgId, router)
	fmt.Println("Add Router success.")
}

func (server *Server) SetConnHook(hook ziface.IConnHook) {
	server.ConnManager.SetConnHook(hook)
}

//初始化server的方法，返回一个抽象层的Server
func NewServer(name string) ziface.IServer {
	if name == "" {
		name = utils.GlobalObject.Name
	}
	s := &Server{
		Name:        name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		ReqToRouter: NewReqHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}
