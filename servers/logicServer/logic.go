package main

import (
	"fmt"
	_ "net/http/pprof"

	"log"
	"net/http"
	"serverFramework/core"
	"serverFramework/internal/utils"
	"serverFramework/internal/version"
)

func main() {
	fmt.Println(version.String("logic Server"))

	core.Run()
	//core.Run("127.0.0.1:8089")
	//core.Run("localhost")
	// core.Run(":8089")

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var wg utils.WaitGroupWrapper
	wg.Wrap(func() {
		serverRoom()
	})

	// sc.
	// sync.WaitGroup.Wait()
}

func serverRoom() {
	for {

	}
}
