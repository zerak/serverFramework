package main

import (
	"fmt"
	_ "net/http/pprof"

	"log"
	"net/http"
	"serverFramework/core"
	"serverFramework/internal/version"
)

func main() {
	fmt.Println(version.String("Chat Server"))

	sc := core.New()
	sc.Run()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// sc.
	// sync.WaitGroup.Wait()
}
