package ziface

import "net"

type IConnection interface {
	// 启动链接 开始当前工kj
	Start()
	// 停止链接 结束当前链接
	Stop()
	// 获取当前链接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	// 获取当前链接的ID
	GetConnID() uint32
	// 获取远程客户端的地址信息
	RemoteAddr() net.Addr
	// 发送消息给客户端
	SendMsg(data []byte) error
}

// 定义处理链接的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
