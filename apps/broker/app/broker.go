package app

import (
	"github.com/flinkcoin/mono/apps/broker/internal/net"
)

type Broker struct {
	Harbor *net.Net
}

func NewBroker(harbor *net.Net) *Broker {
	return &Broker{Harbor: harbor}
}
