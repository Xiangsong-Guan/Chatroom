// Package client 提供了对客户端的实验，并且提供了一些工具
package client

import (
	"Chatroom/hub"
	"Chatroom/message"
	"bytes"
	"fmt"
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
	wsConn *websocket.Conn
	info   message.Info
}

// New generate a Client
func New(hb *hub.Hub, wsConn *websocket.Conn, inform message.Info) (c *Client) {
	c = &Client{
		hub:    hb,
		wsConn: wsConn,
		info:   inform,
		send:   make(chan []byte, 128),
	}
	hb.Register <- c
	msg := message.New([]byte("加入聊天"), time.Now(), inform)
	hb.Boradcast <- msg
	return
}

// ReadPump 函数 ：从websockets 连接的消息发送到到 hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.hub.Boradcast <- message.New([]byte("离开聊天"), time.Now(), c.info)
		c.wsConn.Close()
	}()
	c.wsConn.SetReadLimit(maxMessageSize)
	c.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	c.wsConn.SetPongHandler(func(string) error { c.wsConn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		m := message.New(msg, time.Now(), c.info)
		c.hub.Boradcast <- m
	}
}

// WritePump 函数
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.hub.Unregister <- c
		ticker.Stop()
		c.wsConn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.wsConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SenderChannel 返回信道
func (c *Client) SenderChannel() chan<- []byte {
	return c.send
}

// GenSummary return summary about client
func (c *Client) GenSummary() string {
	return fmt.Sprintln(c.info.Name, c.info.Mail, c.info.IPAddress)
}
