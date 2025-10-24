package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/logger"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name       string
	IPVersion  string
	IP         string
	Port       int
	Log        *logger.Logger
	msgHandler ziface.IMsgHandle
}

func NewServer(name string) ziface.IServer {
	// 初始化全局配置文件
	utils.GlobalObject.Reload()

	return &Server{
		Name:       name,
		IPVersion:  "tcp4",
		IP:         "0.0.0.0",
		Port:       8999,
		Log:        logger.NewLogger(logger.WithGroup("zinx-s")),
		msgHandler: NewMsgHandle(),
	}
}

// ! 这里是一个回显业务， 是server端的业务 目前是固定的 后续需要改成可配置的
func CallBack(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}
func (s *Server) Start() {
	s.Log.Info("server start", "name", s.Name, "ip", s.IP, "port", s.Port)
	s.Log.Info("MetaData", "version", utils.GlobalObject.Version, "max_conn", utils.GlobalObject.MaxConn, "max_packet_size", utils.GlobalObject.MaxPacketSize)
	// 解析ip地址
	go func() {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			s.Log.Error("resolve ip error", "err", err)
			return
		}

		// 监听tcp
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			s.Log.Error("listen error", "err", err)
			return
		}
		s.Log.Info("server start success", "addr", listener.Addr())
		var cid uint32 = 0
		// 阻塞等待
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				s.Log.Error("accept error", "err", err)
				continue
			}
			s.Log.Info("accept success", "remoteAddr", conn.RemoteAddr().String())
			// 创建新的链接对象 并且去调用链接业务
			dealConn := NewConnection(conn, cid, s.msgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	s.Log.Error("server stop not impl")
}

func (s *Server) Serve() {
	s.Start()
	// TODO 可以做其他业务
	select {}
}
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.Log.Info("add router", "msgId", msgId)
	s.msgHandler.AddRouter(msgId, router)
}
