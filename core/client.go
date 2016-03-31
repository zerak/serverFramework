package core

import (
	"bufio"
	"fmt"
	"net"
	"sync"
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

	writeLock sync.RWMutex
	metaLock  sync.RWMutex
}

func newClient(id int64, conn net.Conn, ctx *context) *ClientV1 {
	fmt.Printf("new client ...\n")
	var identifier string
	if conn != nil {
		identifier, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	}

	c := &ClientV1{
		ClientID: identifier,
		HostName: identifier,

		ID:   id,
		Conn: conn,
		ctx:  ctx,

		Reader: bufio.NewReaderSize(conn, defaultBufferSize),
		Writer: bufio.NewWriterSize(conn, defaultBufferSize),
	}
	fmt.Printf("new client [%s]\n",conn.RemoteAddr())
	return c
}
