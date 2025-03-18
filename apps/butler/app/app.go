package app

import "github.com/flinkcoin/mono/apps/butler/internal/messaging"

type App struct {
	Queue *messaging.Queue
}

func NewApp(queue *messaging.Queue) *App {
	return &App{Queue: queue}
}

func (a *App) Connect() {
	a.Queue.Connect()
}
