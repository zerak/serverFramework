package main

import (
	"fmt"
	_ "net/http/pprof"

	"log"
	"net/http"
	"serverFramework/core"
)

func main() {
	fmt.Printf("Hello serverFramework.\n")

	sc := core.New()
	sc.Run()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// sc.
	// sync.WaitGroup.Wait()
}
