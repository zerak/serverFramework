package core

import (
	"bufio"
	"net"
	"sync"
	"time"
)

const defaultBufferSize = 16 * 1024

type ClientV1 struct {
	ClientID string
	HostName string

	net.Conn

	ID  int64
	ctx *context

	// reading/writing interfaces
	Reader *bufio.Reader
	Writer *bufio.Writer

	HeartbeatInterval time.Duration
	MsgTimeout        time.Duration

	ExitChan chan int

	writeLock sync.RWMutex
	metaLock  sync.RWMutex
}

func (c *ClientV1) String() string {
	return c.RemoteAddr().String()
}

func (c *ClientV1) Close() {
	c.Conn.Close()
	close(c.ExitChan)
}

func newClient(id int64, conn net.Conn, ctx *context) *ClientV1 {
	Info("new client ...")
	var identifier string
	if conn != nil {
		identifier, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	}

	c := &ClientV1{
		ClientID: identifier,
		HostName: identifier,

		Conn: conn,

		ID:  id,
		ctx: ctx,

		Reader: bufio.NewReaderSize(conn, defaultBufferSize),
		Writer: bufio.NewWriterSize(conn, defaultBufferSize),

		HeartbeatInterval: 10 * time.Second / 2,
		MsgTimeout:        10 * time.Second / 2,

		ExitChan: make(chan int),
	}
	ServerLogger.Info("new client id[%d] addr[%s] identifier[%s]", c.ID, conn.RemoteAddr(), c.ClientID)
	return c
}
