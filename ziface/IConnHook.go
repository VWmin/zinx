package ziface

type IConnHook interface {
	AfterConnStart(conn IConnection)
	BeforeConnStop(conn IConnection)
}

type BaseConnHook struct {}

func (b *BaseConnHook) AfterConnStart(conn IConnection) {}

func (b *BaseConnHook) BeforeConnStop(conn IConnection) {}
