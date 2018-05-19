// Package hub inplements Hub
package hub

import (
	"Chatroom/historylog"
)

// Hub is a powerful controller for clients
type Hub struct {
	clients    map[Sender]chan<- []byte
	Boradcast  chan Sendable
	Register   chan Sender
	Unregister chan Sender
	history    *historylog.HistoryLog
}

// New a hub
func New(hl *historylog.HistoryLog) (h *Hub) {
	h = &Hub{
		Boradcast:  make(chan Sendable, 256),
		Register:   make(chan Sender),
		Unregister: make(chan Sender),
		clients:    make(map[Sender]chan<- []byte),
		history:    hl,
	}
	go h.run()
	return
}

// Sender 能够提供一个可发送的信道来接受消息
type Sender interface {
	historylog.HistoryLogger
	SenderChannel() chan<- []byte
}

// Sendable 抽象了在系统内部可以被传播的内容
type Sendable interface {
	historylog.HistoryLogger
	Data() []byte
}

func (h *Hub) run() {
	for {
		select {
		case clent := <-h.Register:
			h.clients[clent] = clent.SenderChannel()
		case clent := <-h.Unregister:
			if ch, ok := h.clients[clent]; ok {
				close(ch)
				delete(h.clients, clent)
			}
		case message := <-h.Boradcast:
			for _, ch := range h.clients {
				select {
				case ch <- message.Data():
				default:
					close(ch)
				}
			}
			h.history.PrintFrom(message)
		}
	}
}
