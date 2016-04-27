package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"serverFramework/core"
	"serverFramework/internal/utils"
)

func main() {
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
	wg.Wait()

	// sc.
	// sync.WaitGroup.Wait()
}

func serverRoom() {
	for {

	}
}
