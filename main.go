package main

import (
	"fmt"
	// "sync"

	"serverFramework/core"
)

func main() {
	fmt.Printf("Hello serverFramework.\n")

	sc := core.New()
	sc.Run()

	// sc.
	// sync.WaitGroup.Wait()
}
