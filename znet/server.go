package znet

import (
	"fmt"
	"net"
	"zinx/logger"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Log       *logger.Logger
}

func NewServer(name string) *Server {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Log:       logger.NewLogger(logger.WithGroup("zinx-s")),
	}
}
func (s *Server) Start() {
	s.Log.Info("server start", "name", s.Name, "ip", s.IP, "port", s.Port)
	// 解析ip地址
	go func() {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			s.Log.Error("resolve ip error", "err", err)
			return
		}

		// 监听tcp
		listener, err := net.Listen("tcp", addr.String())
		if err != nil {
			s.Log.Error("listen error", "err", err)
			return
		}
		s.Log.Info("server start success", "addr", listener.Addr())

		// 阻塞等待
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					s.Log.Error("accept error", "err", err)
					continue
				}
				buf := make([]byte, 512)
				cnt, err := conn.Read(buf)
				if err != nil {
					s.Log.Error("read error", "err", err)
					continue
				}
				_, err = conn.Write(buf[:cnt])
				if err != nil {
					s.Log.Error("write error", "err", err)
					continue
				}
			}
		}()
	}()
}

func (s *Server) Stop() {
	panic("not impl")
}

func (s *Server) Serve() {
	s.Start()
	// todo 可以做其他业务
	select {

	}
}
