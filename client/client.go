// Package client 提供了对客户端的实验，并且提供了一些工具
package client

import (
	"Chatroom/hub"
	"Chatroom/message"
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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
func New(hubs *hub.Hub, webConn *websocket.Conn, inform message.Info) (c *Client) {
	c = &Client{
		hub:    hubs,
		webCnn: webConn,
		info:   inform,
		send:   make(chan []byte, 128),
	}
	hubs.Register <- c
	return
}

// ReadPump 函数 ：从websockets 连接的消息发送到到 hub
func (clent *Client) ReadPump() {
	defer func() {
		clent.hub.Unregister <- clent
		clent.webCnn.Close()
	}()
	clent.webCnn.SetReadLimit(maxMessageSize)
	clent.webCnn.SetReadDeadline(time.Now().Add(pongWait))
	clent.webCnn.SetPongHandler(func(string) error { clent.webCnn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := clent.webCnn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		m := message.New(msg, time.Now(), clent.info)
		clent.hub.Boradcast <- m
	}
}

// WritePump 函数
func (clent *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		clent.hub.Unregister <- clent
		ticker.Stop()
		clent.webCnn.Close()
	}()
	for {
		select {
		case msg, ok := <-clent.send:
			clent.webCnn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				clent.webCnn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := clent.webCnn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// Add queued chat messages to the current websocket message.
			n := len(clent.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-clent.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			clent.webCnn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := clent.webCnn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SenderChannel 返回信道
func (clent *Client) SenderChannel() chan<- []byte {
	return clent.send
}
