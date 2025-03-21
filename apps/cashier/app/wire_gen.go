// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/flinkcoin/mono/apps/cashier/internal/messaging"
	"github.com/flinkcoin/mono/apps/cashier/internal/process"
)

// Injectors from wire.go:

func Init() *App {
	processProcess := process.NewProcess()
	queue := messaging.NewQueue(processProcess)
	app := NewApp(queue)
	return app
}
