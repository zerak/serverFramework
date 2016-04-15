package core

import (
	"io"
	"net"

	"serverFramework/internal/protocol"
)

type Handler interface {
	Handle(net.Conn)
}

type ServerHandler struct {
	ctx *context
}

func (sh *ServerHandler) Handle(clientConn net.Conn) {
	sh.ctx.core.log.Info("[ServerHandler::Handle] client[%s]\n", clientConn.RemoteAddr())

	// The client should initialize itself by sending a 4 byte sequence indicating
	// the version of the protocol that it intends to communicate, this will allow us
	// to gracefully upgrade the protocol away from text/line oriented to whatever...
	buf := make([]byte, 4)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		sh.ctx.core.log.Error("[ServerHandler::Handle] ERROR: failed to read protocol version [%s]\n", err)
		return
	}
	protocolMagic := string(buf)
	sh.ctx.core.log.Info("[ServerHandler::Handle] recv[%s]\n", protocolMagic)

	var pro protocol.Protocol
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
		sh.ctx.core.log.Info("[ServerHandler::Handle] ERROR: client[%s] - [%s]\n", clientConn.RemoteAddr(), err)
		return
	}
	sh.ctx.core.log.Info("[ServerHandler::Handle] client exit[%v] - [%v]\n", clientConn.RemoteAddr(), err)
}
