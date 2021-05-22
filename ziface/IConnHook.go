package ziface

type IConnHook interface {
	AfterConnStart(conn IConnection)
	BeforeConnStop(conn IConnection)
}
