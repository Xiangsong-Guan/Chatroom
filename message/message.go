// Package message 实现了对消息的封装，提供了一些对应的工具
package message

import (
	"net"
	"strings"
	"time"
)

// Info 对应一个发送消息的实体，它封装了该实体的信息
type Info struct {
	Name      string
	IPAddress net.Addr
	Mail      string
}

// GenSummary 实现
func (m Message) GenSummary() string {
	return string(m.Msg) + "from" + m.UserInfo.Name + "@" + m.UserInfo.IPAddress.String()
}

// Message 封装了聊天消息，在系统内部，聊天消息以此形式传播
type Message struct {
	SendTime time.Time
	Msg      []byte
	UserInfo Info
}

// New generate a message
func New(message []byte, sendtime time.Time, userinfo Info) Message {
	return Message{
		SendTime: sendtime,
		Msg:      message,
		UserInfo: userinfo,
	}
}

// Data 实现数据流的转换
func (m Message) Data() []byte {
	ips := strings.Split(m.UserInfo.IPAddress.String(), ".")
	return []byte(m.UserInfo.Name + "|" + "*.*." + ips[2] + "." + ips[3] + "|" + m.SendTime.Format(time.RFC1123) + "|" + string(m.Msg))
}
