package core

import (
	"net"
	"sync"
	"time"
)

type Message struct {
	ID        int
	Body      []byte
	Timestamp time.Time
	pool      sync.Pool

	net.Conn
}

func NewMsg(id int, body []byte, client *ClientV1) *Message {
	return &Message{
		ID:        id,
		Body:      body,
		//Timestamp: time.Now().UnixNano(),
		Timestamp: time.Now(),

		Conn: client.Conn,
	}
}

func NewMsgEmpty() *Message {
	return &Message{
		//Timestamp: time.Now().UnixNano(),
		Timestamp: time.Now(),
	}
}
