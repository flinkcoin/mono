//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/flinkcoin/mono/apps/broker/internal/config"
	"github.com/flinkcoin/mono/apps/broker/internal/net"
	"github.com/google/wire"
	"libs/core/pkg"
)

func Init() *Broker {

	wire.Build(pkg.NewLogger, config.NewConfig, net.NewNode, NewBroker)
	return nil

}
