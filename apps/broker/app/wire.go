//go:build wireinject
// +build wireinject

package app

import (
	"github.com/flinkcoin/mono/apps/broker/internal/net"
	"github.com/google/wire"
)

func Init() *Broker {
	wire.Build(net.NewNode, NewBroker)
	return nil
}
