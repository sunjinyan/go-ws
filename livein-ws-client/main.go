package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// WebSocket 服务器地址
	serverAddr := "ws://localhost:8080/ws"

	// 建立连接
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("无法连接到服务器:", err)
	}
	defer conn.Close()
	fmt.Println("成功连接到服务器:", serverAddr)

	// 捕获中断信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 启动一个 Goroutine 用于接收消息
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("读取消息时出错:", err)
				return
			}
			fmt.Printf("收到服务器消息: %s\n", message)
		}
	}()

	// 主线程用于发送消息
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			// 每 2 秒发送一次消息
			message := fmt.Sprintf("客户端时间: %s", t)
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("发送消息时出错:", err)
				return
			}
			fmt.Printf("发送消息: %s\n", message)
		case <-interrupt:
			fmt.Println("收到中断信号，关闭连接...")
			// 发送关闭消息给服务器
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "客户端关闭连接"))
			if err != nil {
				fmt.Println("发送关闭消息时出错:", err)
				return
			}
			// 等待服务器确认关闭
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
