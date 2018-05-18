// Package client 提供了对客户端的实验，并且提供了一些工具
package client

import (
	"Chatroom/hub"
	"Chatroom/message"

	"github.com/gorilla/websocket"
)

// Client 对应着每一个参与聊天的远程实体，它记录了该实体的信息并且负责
// 针对该实体收发消息
type Client struct {
	send   chan []byte
	hub    *hub.Hub
	webCnn *websocket.Conn
	info   message.Info
}

// New generate a Client
func New(hubs *hub.Hub, webConn *websocket.Conn, inform message.Info) Client {
	return Client{
		hub:    hubs,
		webCnn: webConn,
		info:   inform,
	}
}

// readPump 函数
func (client *Client) ReadPump() {
}

// writePump 函数
func (client *Client) WritePump() {
}

// Sender 能够提供一个可发送的信道来接受消息
type Sender interface {
	GetSenderChannel() chan<- []byte
}
