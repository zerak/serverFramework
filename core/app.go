package core

import (
	"net"
	"os"
	"sync"
	"time"

	"serverFramework/internal/moduledb"
	"serverFramework/internal/modulelogic"
	"serverFramework/internal/utils"
)

var (
	// ServerApp is an application instance
	ServerApp *ServerCore
)

type ServerCore struct {
	sync.RWMutex

	clientIDSequence int64     // 64bit atomic vars need to be first for proper alignment on 32bit platforms
	startTime        time.Time // server start time
	tcpListener      net.Listener
	wg               utils.WaitGroupWrapper
	db               chan *moduledb.DBModuler       // the db chan
	logic            chan *modulelogic.LogicModuler // the logic chan
}

func init() {
	ServerApp = New()
}

func New() *ServerCore {
	sc := &ServerCore{
		startTime: time.Now(),
	}

	return sc
}

func (sc *ServerCore) Run() {
	tcpListener, err := net.Listen("tcp", SConfig.TCPAddr)
	if err != nil {
		ServerLogger.Error("listen [%s] failed ->%s", SConfig.TCPAddr, err)
		os.Exit(0)
	}

	sc.Lock()
	sc.tcpListener = tcpListener
	sc.Unlock()

	Info("server listen on", SConfig.TCPAddr)

	ctx := &context{sc}

	handle := &ServerHandler{ctx: ctx}
	// start accept routine
	sc.wg.Wrap(func() {
		HandleAccept(sc.tcpListener, handle)
	})

	// start process msg
	sc.wg.Wrap(func() {

	})

	sc.wg.Wait()
}
