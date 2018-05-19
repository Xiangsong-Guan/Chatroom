// Package hub inplements Hub
package hub

// Hub is a powerful controller for clients
type Hub struct {
	clients    map[*Sender]chan<- []byte
	Boradcast  chan []byte
	Register   chan *Sender
	Unregister chan *Sender
}

// New a hub
func New() *Hub {
	return &Hub{
		Boradcast:  make(chan []byte, 256),
		Register:   make(chan *Sender),
		Unregister: make(chan *Sender),
		clients:    make(map[*Sender]chan<- []byte),
	}
}

// Sender 能够提供一个可发送的信道来接受消息
type Sender interface {
	SenderChannel() chan<- []byte
}

// Run  hub 的业务逻辑处理
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = (*client).SenderChannel()
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients)
			}
		}
	}
}
