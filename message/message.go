// Package message 实现了对消息的封装，提供了一些对应的工具
package message

import (
	"strings"
	"time"
)

// SPChar pre defined char
var SPChar = "\t"

// Info 对应一个发送消息的实体，它封装了该实体的信息
type Info struct {
	Name      string
	IPAddress string
	Mail      string
}

// GenSummary 实现
func (m Message) GenSummary() string {
	return string(m.msg) + " from " + m.senderInfo.Name + "@" + m.senderInfo.IPAddress
}

// Message 封装了聊天消息，在系统内部，聊天消息以此形式传播
type Message struct {
	sendTime   time.Time
	msg        []byte
	senderInfo Info
}

// New generate a message
func New(message []byte, sendtime time.Time, userinfo Info) Message {
	return Message{
		sendTime:   sendtime,
		msg:        message,
		senderInfo: userinfo,
	}
}

// Data 实现数据流的转换
func (m Message) Data() []byte {
	ips := strings.Split(m.senderInfo.IPAddress, ".")
	return []byte(m.senderInfo.Name + SPChar + ips[2] + "." + ips[3] + SPChar + m.sendTime.Format(time.RFC1123) + SPChar + string(m.msg))
}
