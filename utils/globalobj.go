package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         `json:"host"`     // TCP服务地址
	TcpPort   int            `json:"tcp_port"` // TCP服务端口号
	Name      string         `json:"name"`     // 服务器名称

	Version          string `json:"version"`             // Zinx版本号
	MaxPacketSize    uint32 `json:"mac_packet_size"`     // 数据包的最大值
	MaxConn          uint32 `json:"mac_conn"`            // 当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 `json:"worker_pool_size"`    // 工作池大小
	MaxWorkerTaskLen uint32 `json:"max_worker_task_len"` // 当前工作池最大任务长度
	MaxMsgChanLen    uint32 `json:"max_msg_chan_len"`    // 最大的消息通道数量
	// config
	ConfFilePath string `json:"conf_file_path"`
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.8",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxPacketSize:    4096,
		MaxConn:          12000,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "conf/zinx.json",
		MaxMsgChanLen:    1024,
	}
	GlobalObject.Reload()
}
