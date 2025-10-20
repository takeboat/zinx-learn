package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

// mock 模拟客户端
func main() {
	fmt.Println("client start...")
	// 创建链接 获取conn
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
	for {
		// 发送数据
		_, err := conn.Write([]byte("Hello from client\n"))
		if err != nil {
			slog.Error("write error:", "err", err)
			return
		}
		// 读取响应
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("read error:", "err", err)
			return
		}
		slog.Info("receive from server:", "data", string(buf[:n]))
		// cpu 休眠1秒
		time.Sleep(1 * time.Second)
	}
}
