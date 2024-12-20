package internal

import (
	"github.com/flinkcoin/mono/apps/broker/internal/config"
	"github.com/flinkcoin/mono/apps/broker/internal/net"
	"log/slog"
)

var logger *slog.Logger
var conf *config.Config

type Broker struct {
	Harbor *net.Net
}

func NewBroker(l *slog.Logger, c *config.Config, harbor *net.Net) *Broker {
	conf = c
	logger = l

	return &Broker{Harbor: harbor}
}
