// Package message 实现了对消息的封装，提供了一些对应的工具
package message

import "time"

// Info 对应一个发送消息的实体，它封装了该实体的信息
type Info struct {
	Name      string
	IPAddress string
	Mail      string
}

// GenSummary 实现
func (m Message) GenSummary() string {
	return string(m.Message) + "from" + m.UserInfo.Name + "@" + m.UserInfo.IPAddress
}

// Message 封装了聊天消息，在系统内部，聊天消息以此形式传播
type Message struct {
	SendTime time.Time
	Message  []byte
	UserInfo Info
}

// New generate a message
func New(message []byte, sendtime time.Time, userinfo Info) Message {
	return Message{
		SendTime: sendtime,
		Message:  message,
		UserInfo: userinfo,
	}
}

// Sendable 抽象了在系统内部可以被传播的内容
type Sendable interface {
	Data() []byte
}
