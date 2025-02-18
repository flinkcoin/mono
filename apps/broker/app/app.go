package app

import (
	"github.com/flinkcoin/mono/apps/broker/internal/networking"
)

type App struct {
	Host *networking.Host
}

func NewApp(host *networking.Host) *App {
	return &App{Host: host}
}
