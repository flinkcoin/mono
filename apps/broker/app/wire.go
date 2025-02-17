//go:build wireinject
// +build wireinject

package app

import (
	"github.com/flinkcoin/mono/apps/broker/internal/networking"
	"github.com/google/wire"
)

func Init() *Broker {
	wire.Build(networking.NewNode, NewBroker)
	return nil
}
