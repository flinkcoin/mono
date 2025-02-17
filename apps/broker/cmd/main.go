package main

import (
	"fmt"
	"github.com/flinkcoin/mono/apps/broker/internal"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	/*broker, err :=*/
	app := internal.Init()
	app.Net.Init()

	if len(os.Args) >= 2 {
		fmt.Println("Usage: program <argument>")
		app.Net.Connect(os.Args[1])
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	//	host.Init()
}
