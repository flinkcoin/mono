package main

import (
	"fmt"
	"github.com/flinkcoin/mono/apps/broker/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	/*broker, err :=*/
	a := app.Init()
	a.Host.Init()

	if len(os.Args) >= 2 {
		fmt.Println("Usage: program <argument>")
		a.Host.Connect(os.Args[1])
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	//	host.Init()
}
