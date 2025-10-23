package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"
	"zinx/znet"
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
	var id uint32
	dp := znet.NewDataPack()
	for {
		// 发送数据
		id++
		msg := znet.NewMsgPackage(id, []byte("Hello from client"))
		data, err := dp.Pack(msg)
		if err != nil {
			fmt.Println("pack error:", err)
			return
		}
		_, err = conn.Write(data)
		if err != nil {
			slog.Error("write error:", "err", err)
			return
		}
		slog.Info("write success...")
		// cpu 休眠1秒
		time.Sleep(1 * time.Second)
	}
}
