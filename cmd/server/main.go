package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"zinx/ziface"
	"zinx/znet"
)

func main() {
	// server := znet.NewServer("[Zinx]")
	// server.AddRouter(&PingRouter{})
	// server.Serve()
	listener, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			continue
		}
		go func(conn net.Conn) {
			dp := znet.NewDataPack()
			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error:", err)
					break
				}
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("unpack err:", err)
					continue
				}
				if msgHead.GetDataLen() > 0 {
					data := make([]byte, msgHead.GetDataLen())
					_, err := io.ReadFull(conn, data)
					if err != nil {
						fmt.Println("read msg data error:", err)
						break
					}
					slog.Info("receive msg:", "msgID:", msgHead.GetMsgId(), "dataLen:", msgHead.GetDataLen(), "data:", string(data))
				}
			}
		}(conn)
	}
}

type PingRouter struct {
	znet.BaseRouter
}

func (r *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call PingRouter PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before-----ping...ping...ping\n"))
	if err != nil {
		fmt.Println("CallBackToClient error")
	}
}

func (r *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Handle-----ping...ping...ping\n"))
	if err != nil {
		fmt.Println("CallBackToClient error")
	}
}
func (r *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call PingRouter PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("post-----ping...ping...ping\n"))
	if err != nil {
		fmt.Println("CallBackToClient error")
	}
}
