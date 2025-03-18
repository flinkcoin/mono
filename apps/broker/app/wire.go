//go:build wireinject
// +build wireinject

package app

import (
	"github.com/flinkcoin/mono/apps/broker/internal/messaging"
	"github.com/flinkcoin/mono/apps/broker/internal/networking"
	"github.com/google/wire"
)

func Init() *App {
	wire.Build(networking.NewHost, messaging.NewQueue, NewApp)
	return nil
}
