package test

import (
	"fmt"
	"github.com/vwmin/zinx/znet"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/
	// 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		return
	}

	// 从客户端读数据，拆包处理
	go func() {
		for true {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			go func(conn net.Conn) {
				dp := znet.NewDataPack()
				for true {
					// 从conn中读出head
					headBytes := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headBytes); err != nil {
						return
					}

					head, err := dp.Unpack(headBytes)
					if err != nil {
						return
					}
					if head.GetDataLen() > 0 {
						// 有数据，从conn中读出data
						data := make([]byte, head.GetDataLen())
						if _, err := io.ReadFull(conn, data); err != nil {
							return
						}
						head.SetData(data)
					}

					fmt.Printf("id:%d, len:%d, content:%s\n", head.GetMsgId(), head.GetDataLen(), head.GetData())
				}
			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		return
	}

	dp := znet.NewDataPack()

	msg1 := znet.NewMessage(1, []byte("hello"))
	msg2 := znet.NewMessage(2, []byte("world"))

	toSend1, _ := dp.Pack(msg1)
	toSend2, _ := dp.Pack(msg2)

	toSend1 = append(toSend1, toSend2...)

	conn.Write(toSend1)

	select {}

}
