package main

import (
	"github.com/flinkcoin/mono/apps/butler/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	/*broker, err :=*/
	a := app.Init()
	a.Connect()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
