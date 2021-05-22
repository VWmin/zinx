package utils

import (
	"encoding/json"
	"fmt"
	"github.com/vwmin/zinx/ziface"
	"io/ioutil"
)

/**
  存储框架全局参数，供其他模块使用
  一些参数可通过zinx.json配置
*/
type GlobalObj struct {

	/*
		Server
	*/

	// 当前Zinx全局的Server对象
	TcpServer ziface.IServer

	// 当前服务器主机监听的IP
	Host string

	// 当前服务器主机监听的Port
	TcpPort int

	// 当前服务器名称
	Name string

	// 最大连接个数
	MaxCoon int

	// 数据包最大值
	MaxPackageSize uint32

	// 工作线程数量
	WorkerSize uint

	// 允许的最大工作线程数量
	maxWorkerSize uint

	// 每个worker的任务队列大小
	WorkerTaskQueueSize uint

	// 最大的任务队列大小
	maxWorkerTaskQueueSize uint

	/*
		Zinx
	*/

	// 版本号
	Version string
}

// 定义一个全局的对外对象
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	fmt.Println("Trying resolve config file...")
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("Loading config file failed, absorbing...")
		return
	}

	// 将json内容解析到GlobalObject中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		fmt.Println("Resolving config file failed, exiting...")
		panic(err)
	}

	// 最大值检查
	if g.WorkerSize > g.maxWorkerSize {
		g.WorkerSize = g.maxWorkerSize
	}

	if g.WorkerSize > g.maxWorkerSize {
		g.WorkerSize = g.maxWorkerSize
	}

}

// 提供init()方法，初始化全局对象
func init() {
	// 无配置文件时默认值
	GlobalObject = &GlobalObj{
		Host:                   "0.0.0.0",
		TcpPort:                8999,
		Name:                   "ZinxServerApp",
		MaxCoon:                1000,
		MaxPackageSize:         4096,
		WorkerSize:             10,
		maxWorkerSize:          1024,
		WorkerTaskQueueSize:    20,
		maxWorkerTaskQueueSize: 1024,
		//Version:        "V0.5",
	}

	// 尝试读取配置文件
	GlobalObject.Reload()
}
