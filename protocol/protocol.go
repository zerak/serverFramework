package protocol

import (
	"bufio"
	"net"
	"strings"
)

type Protocoler interface {
	IOLoop(conn net.Conn) error
	Send(w *bufio.Writer, data []byte) (int, error)
}

var proAdapts = make(map[string]Protocoler)

func Register(name string, pro Protocoler) {
	if pro == nil {
		panic("protocol: Register pro is nil")
	}

	if _, ok := proAdapts[strings.ToUpper(name)]; ok {
		panic("protocol: Register called twice for adapter " + name)
	}
	proAdapts[name] = pro
}
