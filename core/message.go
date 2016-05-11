package core

import (
	"net"
	"sync"
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
	pool      sync.Pool

	net.Conn
}

func NewMsg(id int, body []byte, client *ClientV1) *Message {
	return &Message{
		ID:        id,
		Body:      body,
		Timestamp: time.Now().UnixNano(),

		Conn: client.Conn,
	}
}

func NewMsgEmpty() *Message {
	return &Message{
		Timestamp: time.Now().UnixNano(),
	}
}
