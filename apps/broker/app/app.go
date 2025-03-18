package app

import (
	"github.com/flinkcoin/mono/apps/broker/internal/messaging"
	"github.com/flinkcoin/mono/apps/broker/internal/networking"
)

type App struct {
	Host  *networking.Host
	Queue *messaging.Queue
}

func NewApp(host *networking.Host, queue *messaging.Queue) *App {
	return &App{
		Host:  host,
		Queue: queue,
	}
}
