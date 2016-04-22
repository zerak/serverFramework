package core

import (
	"net"
	"time"
)

const (
	MsgIDLength       = 16
	minValidMsgLength = MsgIDLength + 8 + 2 // Timestamp + Attempts
)

type MessageID [MsgIDLength]byte

type Message struct {
	ID        int
	Body      []byte
	Timestamp int64
	Attempts  uint16

	net.Conn
}

func NewMessage(id int, body []byte, client *ClientV1) *Message {
	return &Message{
		ID:        id,
		Body:      body,
		Timestamp: time.Now().UnixNano(),

		Conn: client.Conn,
	}
}
