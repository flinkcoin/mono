package app

import (
	"github.com/flinkcoin/mono/apps/broker/internal/networking"
)

type Broker struct {
	Harbor *networking.Host
}

func NewBroker(harbor *networking.Host) *Broker {
	return &Broker{Harbor: harbor}
}
