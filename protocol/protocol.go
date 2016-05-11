package protocol

import (
	"net"
	"strings"

	. "serverFramework/client"
)

// Protocol describes the basic behavior of any protocol in the system
type Protocol interface {
	IOLoop(conn net.Conn) error
	Send(c Client, data []byte) (int, error)
}

var ProAdapts = make(map[string]Protocol)

func Register(name string, pro Protocol) {
	if pro == nil {
		panic("protocol: Register pro is nil")
	}

	if _, ok := ProAdapts[strings.ToUpper(name)]; ok {
		panic("protocol: Register called twice for adapter " + name)
	}
	ProAdapts[strings.ToUpper(name)] = pro
}
