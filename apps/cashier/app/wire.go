//go:build wireinject
// +build wireinject

package app

import (
	"github.com/flinkcoin/mono/apps/broker/internal/networking"
	"github.com/google/wire"
)

func Init() *App {
	wire.Build(networking.NewHost, NewApp)
	return nil
}
