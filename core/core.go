package core

import (
	"strconv"
	"strings"
)

const (
	// VERSION represent beego web framework version.
	VERSION = "1.6.1"

	// DEV is for develop
	DEV = "dev"
	// PROD is for production
	PROD = "prod"
)

func initBeforeRun() {
}

// Run application.
// core.Run() default run on ListenPort
// core.Run("localhost")
// core.Run(":8089")
// core.Run("127.0.0.1:8089")
func Run(params ...string) {
	Info("server run ...")

	initBeforeRun()

	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			addr := strs[0]
			BeeLogger.Info("listen on addr [%v]", addr)
		}
		if len(strs) > 1 && strs[1] != "" {
			port, _ := strconv.Atoi(strs[1])
			BeeLogger.Info("listen on port [%v]", port)
		}
	}

	Server.Run()

	Info("server exit")
}
