package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

type GlobalObj struct {
	TcpServer     ziface.IServer // 当前Zinx全局的Server对象
	Host          string         `json:"host"`            // TCP服务地址
	TcpPort       int            `json:"tcp_port"`        // TCP服务端口号
	Name          string         `json:"name"`            // 服务器名称
	Version       string         `json:"version"`         // Zinx版本号
	MaxPacketSize uint32         `json:"mac_packet_size"` // 数据包的最大值
	MaxConn       uint32         `json:"mac_conn"`        // 当前服务器主机允许的最大链接个数
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
		Name:          "ZinxServerApp",
		Version:       "V0.4",
		TcpPort:       8999,
		Host:          "0.0.0.0",
		MaxPacketSize: 4096,
		MaxConn:       12000,
	}
	GlobalObject.Reload()
}
