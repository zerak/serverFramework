package core

import (
	"net"
	"runtime"
	"strings"
)

type Acceptor struct {
}

func HandleAccept(listener net.Listener, handler Handler) {
	Info("start handle accept ...")

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				Error("NOTICE: temporary Accept() failure -", err)
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				Error("ERROR: listener.Accept() -", err)
			}
			break
		}

		Info("new client", clientConn.RemoteAddr())

		go handler.Handle(clientConn)
	}

	Error("accept routine exit")
}
