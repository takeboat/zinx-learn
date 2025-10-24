package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	fmt.Println("client link success...")
	PingId := uint32(1)
	helloId := uint32(2)
	dp := znet.NewDataPack()
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-exitChan:
				return
			default:
			}
			// 发送数据
			msg := znet.NewMsgPackage(PingId, []byte("Hello from client"))
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
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-exitChan:
				return
			default:
			}
			msg := znet.NewMsgPackage(helloId, []byte("Hello from client"))
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
	}()

	<-exitChan
	wg.Wait()
	fmt.Println("client exit...")
}
