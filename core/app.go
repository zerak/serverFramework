package core

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"serverFramework/utils"
)

const (
	// VERSION represent server framework version.
	VERSION = "0.0.1"

	// DEV is for develop
	DEV = "dev"
	// PROD is for production
	PROD = "prod"
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
	addr := SConfig.TCPAddr
	if SConfig.TCPPort != 0 {
		addr = fmt.Sprintf("%s:%d", SConfig.TCPAddr, SConfig.TCPPort)
	}

	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		ServerLogger.Error("listen [%s:%s] failed ->%s", SConfig.TCPAddr, SConfig.TCPPort, err)
		os.Exit(0)
	}

	sc.Lock()
	sc.tcpListener = tcpListener
	sc.Unlock()

	Info("server listen on", addr)

	handle := &AcceptorHandler{}
	// start accept routine
	sc.wg.Wrap(func() {
		HandleAccept(sc.tcpListener, handle)
	})

	// start process msg ??
	sc.wg.Wrap(func() {

	})

	sc.wg.Wait()
}

func (sc *ServerCore) Version(app string) {
	ServerLogger.Info("%s BASED ON SF v%s(built w/%s)", app, VERSION, runtime.Version())
}
