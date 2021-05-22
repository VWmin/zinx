package znet

import (
	"github.com/vwmin/zinx/ziface"
)

type BaseConnHook struct {}

func (b *BaseConnHook) AfterConnStart(conn ziface.IConnection) {}

func (b *BaseConnHook) BeforeConnStop(conn ziface.IConnection) {}
