package main

import (
	"encoding/pem"
	"fmt"
	"github.com/flinkcoin/mono/apps/broker/app"
	"github.com/flinkcoin/mono/libs/shared/pkg/base"
	"github.com/flinkcoin/mono/libs/shared/pkg/topics"
	"github.com/libp2p/go-libp2p/core/crypto"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	a := app.Init()
	a.Queue.Connect()

	a.Queue.Publish(topics.CashierInbound.String(), []byte("data1"))
	a.Queue.Publish(topics.CashierInbound.String(), []byte("data2"))
	a.Queue.Publish(topics.CashierInbound.String(), []byte("data3"))
	//a.Queue.Subscribe()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	//	host.Init()
}
