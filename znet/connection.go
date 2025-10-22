package znet

import (
	"net"
	"zinx/logger"
	"zinx/ziface"
)

/*
链接模块
*/
type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	// 当前链接绑定的处理业务的方法
	handleApi ziface.HandleFunc
	// 退出的channel
	ExitChan chan struct{}
	// 增加日志
	Log *logger.Logger
}

func NewConnection(conn *net.TCPConn, connID uint32, callback ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleApi: callback,
		isClosed:  false,
		ExitChan:  make(chan struct{}),
		Log:       logger.NewLogger(logger.WithGroup("connection")),
	}
	return c
}

func (c *Connection) StartReader() {
	c.Log.Info("reader goroutine is running", "connID", c.ConnID)

	defer c.Log.Info("reader goroutine is exit", "connID", c.ConnID, "remoteAddr", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端数据到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			c.Log.Error("read from client failed", "connID", c.ConnID, "remoteAddr", c.RemoteAddr().String(), "error", err)
			continue
		}
		// 调用当前链接的业务方法
		if err := c.handleApi(c.Conn, buf, cnt); err != nil {
			c.Log.Error("handle msg failed", "connID", c.ConnID, "remoteAddr", c.RemoteAddr().String(), "error", err)
			close(c.ExitChan)
			return
		}
	}
}
func (c *Connection) Start() {
	c.Log.Info("connection start", "connID", c.ConnID)
	go c.StartReader()
	<-c.ExitChan
	c.Log.Info("connection stop", "connID", c.ConnID)
}
func (c *Connection) Stop() {
	c.Log.Info("connection stop", "connID", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	// 关闭socket链接
	c.Conn.Close()
	// 关闭退出channel
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Connection) SendMsg(data []byte) error {
	return nil
}
