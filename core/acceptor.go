package core

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

type Acceptor struct {
}

func HandleAccept(listener net.Listener, handler Handler) {
	fmt.Printf("[Acceptor::HandleAccept]...\n")

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				fmt.Printf("[Acceptor::HandleAccept] NOTICE: temporary Accept() failure - %s\n", err)
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Printf("[Acceptor::HandleAccept] ERROR: listener.Accept() - %s\n", err)
			}
			break
		}

		fmt.Printf("[Acceptor::HandleAccept] new client(%s)\n", clientConn.RemoteAddr())

		go handler.Handle(clientConn)
	}

	fmt.Printf("[Acceptor::HandleAccept] routines exit\n")
}
