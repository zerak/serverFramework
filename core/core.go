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
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	clientIDSequence int64
	sync.RWMutex
	startTime   time.Time
	tcpListener net.Listener
	wg          util.WaitGroupWrapper

	db chan *DBModuler

}

func (sc *ServerCore) Run() {
	fmt.Printf("[core] Run ...\n")

	tcpListener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Printf("[ServerCore:Run]FATAL: listen (%s) failed - %s\n", "localhost", err)
		os.Exit(0)
	}

	sc.Lock()
	sc.tcpListener = tcpListener
	sc.Unlock()

	fmt.Printf("[ServerCore:Run] server listen on 8888\n")

	ctx := &context{sc}

	handle := &ServerHandler{ctx: ctx}
	// start accept routine
	sc.wg.Wrap(func() {
		HandleAccept(sc.tcpListener, handle)
	})

	fmt.Printf("[core] Run\n")

	sc.wg.Wait()
}

func New() *ServerCore {
	sc := &ServerCore{
		startTime: time.Now()}

	return sc
}
