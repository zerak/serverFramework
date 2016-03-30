package core

import (
	// "fmt"
	"net"
)

type Handler interface {
	Handle(net.Conn)
}

func Handle(clientConn net.Conn) {

}