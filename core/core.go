package core

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"serverFramework/internal/util"
)

type ServerCore struct {
	startTime time.Time
	sync.RWMutex
	tcpListener net.Listener
	wg          util.WaitGroupWrapper
}

func (sc *ServerCore) Init() {
	fmt.Printf("[core] Main ...\n")

	tcpListener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Printf("[ServerCore:Init]FATAL: listen (%s) failed - %s\n", "localhost", err)
		os.Exit(0)
	}

	sc.Lock()
	sc.tcpListener = tcpListener
	sc.Unlock()

	ctx := &context(sc)
	// start goroutine
	sc.wg.Wrap(func() {
		HandleAccept(sc.tcpListener,Hanlde())
	})

	fmt.Printf("[core] Main\n")
}

func New() *ServerCore {
	sc := &ServerCore{
		startTime: time.Now()
	}

	return sc
}
