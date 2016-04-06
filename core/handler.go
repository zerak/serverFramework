package core

import (
	"fmt"
	"io"
	"net"
)

type Handler interface {
	Handle(net.Conn)
}

type ServerHandler struct {
	ctx *context
}

func (sh *ServerHandler) Handle(clientConn net.Conn) {
	fmt.Printf("[ServerHandler::Handle] client[%s]\n", clientConn.RemoteAddr())

	// The client should initialize itself by sending a 4 byte sequence indicating
	// the version of the protocol that it intends to communicate, this will allow us
	// to gracefully upgrade the protocol away from text/line oriented to whatever...
	buf := make([]byte, 4)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		fmt.Printf("[ServerHandler::Handle] ERROR: failed to read protocol version [%s]\n", err)
		return
	}
	protocolMagic := string(buf)
	fmt.Printf("[ServerHandler::Handle] recv[%s]\n", protocolMagic)

	var pro Protocol
	switch protocolMagic {
	case "json":
		pro = &ProtocolJson{ctx: sh.ctx}
	case "  V1", "  v1":
		pro = &ProtocolV1{ctx: sh.ctx}
	case "  V2", "  v2":
		pro = &ProtocolV1{ctx: sh.ctx}
	default:
		return
	}

	err = pro.IOLoop(clientConn)
	if err != nil {
		fmt.Printf("[ServerHandler::Handle] ERROR: client[%s] - [%s]\n", clientConn.RemoteAddr(), err)
		return
	}
	fmt.Printf("[ServerHandler::Handle] client exit[%v] - [%v]\n", clientConn.RemoteAddr(), err)
}
