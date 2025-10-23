package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

func main() {
	server := znet.NewServer("[Zinx]")
	server.AddRouter(&PingRouter{})
	server.Serve()
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
