package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket 升级器配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有连接源（实际使用中请添加安全限制）
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级为 WebSocket 失败:", err)
		return
	}
	defer conn.Close()

	fmt.Println("新的 WebSocket 连接已建立")

	for {
		// 读取消息
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("读取消息时出错:", err)
			break
		}

		fmt.Printf("收到消息: %s\n", message)

		// 发送消息回客户端（回显）
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("发送消息时出错:", err)
			break
		}
	}

}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	port := "8080"
	fmt.Printf("WebSocket 服务器启动，监听端口 %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
