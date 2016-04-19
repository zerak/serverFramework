package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"serverFramework/core"
)

func main() {
	sc := core.New()
	sc.Run()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// sc.
	// sync.WaitGroup.Wait()
}
