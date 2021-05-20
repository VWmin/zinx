package znet

import "github.com/vwmin/zinx/ziface"

type Request struct {

	// 已经和client建立好的连接
	conn ziface.IConnection

	// 客户端请求的数据
	data []byte
}

// 得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 得到请求数据
func (r *Request) GetRequestData() []byte {
	return r.data
}


