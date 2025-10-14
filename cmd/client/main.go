package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("client start...")
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("dial error:", err)
		panic(err)
	}
	defer conn.Close()

	defer func() {
		fmt.Println("client exit...")
	}()
	fmt.Println("client link success...")

	// 发送数据
	conn.Write([]byte("Hello from client\n"))

	// 读取响应
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err == nil {
		fmt.Printf("收到响应: %s\n", string(buf[:n]))
	}
}
