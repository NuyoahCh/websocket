// Package main
// @Author NuyoahCh
// @Date 2025/2/12 23:13
// @Desc
package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients. 定义
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// newHub只是新建了一个空白的Hub，初始化
func newHub() *Hub {
	return &Hub{
		// 用于发送广播的channel
		broadcast: make(chan []byte), // 把消息存到这个channel后，之后会有其它goroutine遍历clients，把消息发送给所有客户端
		// 用于注册客户端的channel
		register: make(chan *Client), // 每当有客户端建立websocket连接时，通过register，把客户端保存到clients引用中
		// 用于注销客户端的channel
		unregister: make(chan *Client), // 每当有客户端断开websocket连接时，通过unregister，把客户端引用从clients中删除
		// 保存了每个客户端的引用的Map
		clients: make(map[*Client]bool), // 其实这个Map的value没有用到，key是客户端的引用，可以当作是其它语言的set
	}
}

// 服务开启时启动的goroutine: hub.run()
func (h *Hub) run() {
	// 创建死循环，不断的从 channel 中读取新的数据
	for {
		// select case 的选择判断方式
		select {
		// 读取到register，就注册客户端
		case client := <-h.register:
			h.clients[client] = true
		// 读取到unregister，就断开客户端连接，删除引用
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// 读取到broadcast，就遍历clients，广播消息（通过把消息写入每个客户端的client.sendchannel 中，实现广播）
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
