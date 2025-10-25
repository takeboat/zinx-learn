package znet

import (
	"errors"
	"io"
	"net"
	"zinx/logger"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	// 退出的channel
	ExitChan chan struct{}
	// 增加日志
	Log        *logger.Logger
	MsgHandler ziface.IMsgHandle
	// 无缓冲管道， 用于读写两个goroutine之间的消息通信
	msgChan     chan []byte
	msgBuffChan chan []byte
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		ExitChan:    make(chan struct{}),
		Log:         logger.NewLogger(logger.WithGroup("connection")),
		MsgHandler:  msgHandler,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}
	// conn added to server mgr
	c.TcpServer.GetConnMgr().Add(c)
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}
func (c *Connection) StartWriter() {
	c.Log.Info("writer goroutine is running")
	defer c.Log.Info("writer goroutine is exit")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				c.Log.Error("write msg error:", "err", err)
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					c.Log.Error("write msg error, conn writer exit", "error", err)
					return
				}
			}
			c.Log.Error("msgBuffChan is closed")
			return
		case <-c.ExitChan:
			return
		}
	}
}
func (c *Connection) Start() {
	c.Log.Info("connection start", "connID", c.ConnID)
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallOnConnStart(c)
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
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	// 关闭退出channel
	c.TcpServer.GetConnMgr().Remove(c) // 删除conn连接 从connmgr中
	close(c.ExitChan)
	close(c.msgBuffChan)
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
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		c.Log.Error("pack error:", "err", err)
		return err
	}
	c.msgChan <- msg
	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		c.Log.Error("pack error:", "err", err)
		return err
	}
	c.msgBuffChan <- msg
	return nil
}
