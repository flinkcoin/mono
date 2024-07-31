//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"libs/core"
	"log/slog"
)

func NewLogger() *slog.Logger {
	wire.Build(core.CreateLogger)
	return &slog.Logger{}
}

func NewConfig() *Config {
	wire.Build(CreateConfig)
	return &Config{}
}
