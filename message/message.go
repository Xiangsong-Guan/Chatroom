// Package message 实现了对消息的封装，提供了一些对应的工具
package message

// Info 对应一个发送消息的实体，它封装了该实体的信息
type Info struct {
}

// Message 封装了聊天消息，在系统内部，聊天消息以此形式传播
type Message struct {
}

// Sendable 抽象了在系统内部可以被传播的内容
type Sendable interface {
	Data() []byte
}
