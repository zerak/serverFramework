package core

import (
	"io"
	"net"
	"strings"

	"serverFramework/protocol"
)

type Handler interface {
	Handle(net.Conn)
}

type AcceptorHandler struct {
}

func (sh *AcceptorHandler) Handle(clientConn net.Conn) {
	ServerLogger.Info("handle client[%s]", clientConn.RemoteAddr())

	// The client should initialize itself by sending a 4 byte sequence indicating
	// the version of the protocol that it intends to communicate, this will allow us
	// to gracefully upgrade the protocol away from text/line oriented to whatever...
	buf := make([]byte, 4)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		ServerLogger.Error("failed to read protocol version ->%s", err)
		return
	}
	protocolMagic := string(buf)
	protocolMagic = strings.ToUpper(protocolMagic)
	ServerLogger.Info("recv client [%v] protocol[%s]", clientConn.RemoteAddr(), protocolMagic)

	if _, ok := protocol.ProAdapts[protocolMagic]; ok {
		err = protocol.ProAdapts[protocolMagic].IOLoop(clientConn)
		if err != nil {
			ServerLogger.Error("client[%s] - [%s]", clientConn.RemoteAddr(), err)
			return
		}
	} else {
		ServerLogger.Warn("the protocol [%v] not support", protocolMagic)
		return
	}

	ServerLogger.Warn("client exit[%v] - [%v]", clientConn.RemoteAddr(), err)
}
