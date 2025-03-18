//go:build wireinject
// +build wireinject

package app

import (
	"github.com/flinkcoin/mono/apps/butler/internal/messaging"
	"github.com/flinkcoin/mono/apps/butler/internal/process"
	"github.com/google/wire"
)

func Init() *App {
	wire.Build(process.NewProcess, messaging.NewQueue, NewApp)
	return nil
}
