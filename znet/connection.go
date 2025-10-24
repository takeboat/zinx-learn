package znet

import (
	"io"
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
	// 退出的channel
	ExitChan chan struct{}
	// 增加日志
	Log        *logger.Logger
	MsgHandler ziface.IMsgHandle
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan struct{}),
		Log:        logger.NewLogger(logger.WithGroup("connection")),
		MsgHandler: msgHandler,
	}
	return c
}

func (c *Connection) StartReader() {
	c.Log.Info("reader goroutine is running", "connID", c.ConnID)

	defer c.Log.Info("reader goroutine is exit", "connID", c.ConnID, "remoteAddr", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端数据到buf中
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			c.Log.Error("read head error:", "err", err)
			break
		}
		// 拆包
		msg, err := dp.Unpack(headData)
		if err != nil {
			c.Log.Error("unpack error:", "err", err)
			break
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				c.Log.Error("read msg data error:", "err", err)
				continue
			}
		}
		msg.SetData(data)
		req := &Request{conn: c, msg: msg}
		go c.MsgHandler.DoMsgHandler(req)
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
