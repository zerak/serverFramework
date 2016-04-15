package core

import (
	"net"
	"os"
	"sync"
	"time"
	//"path/filepath"

	"github.com/astaxie/beego/logs"

	"serverFramework/internal/moduledb"
	"serverFramework/internal/modulelogic"
	"serverFramework/internal/util"
)

type ServerCore struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	clientIDSequence int64
	sync.RWMutex
	startTime   time.Time
	tcpListener net.Listener
	wg          util.WaitGroupWrapper

	db    chan *moduledb.DBModuler
	logic chan *modulelogic.LogicModuler
	log 	*logs.BeeLogger
}

func (sc *ServerCore) Run() {
	defer sc.log.Flush()

	sc.log.EnableFuncCallDepth(true)
	//sc.log.SetLogFuncCallDepth(3)

	sc.log.SetLogger("console", "")
	sc.log.SetLogger("file",`{"filename":"blog.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

	//AppPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//os.Chdir(AppPath)
	//appConfigPath := filepath.Join(AppPath, "conf", "seelog.xml")

	sc.log.Info("Run ...")

	tcpListener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		sc.log.Error("[ServerCore:Run]FATAL: listen (%s) failed - %s", "localhost", err)
		os.Exit(0)
	}

	sc.Lock()
	sc.tcpListener = tcpListener
	sc.Unlock()

	sc.log.Info("[ServerCore:Run] server listen on 8888")

	ctx := &context{sc}

	handle := &ServerHandler{ctx: ctx}
	// start accept routine
	sc.wg.Wrap(func() {
		HandleAccept(sc.tcpListener, handle)
	})

	sc.log.Info("[core] Run")

	sc.wg.Wait()
}

func New() *ServerCore {
	sc := &ServerCore{
		startTime: time.Now(),
		log:logs.NewLogger(100),
	}

	return sc
}
