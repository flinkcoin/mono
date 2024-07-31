package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Config struct {
	Home         string         `env:"HOME"`
	Port         int            `env:"PORT" envDefault:"3000"`
	Password     string         `env:"PASSWORD,unset"`
	IsProduction bool           `env:"PRODUCTION"`
	Duration     time.Duration  `env:"DURATION"`
	Hosts        []string       `env:"HOSTS" envSeparator:":"`
	TempFolder   string         `env:"TEMP_FOLDER,expand" envDefault:"${HOME}/tmp"`
	StringInts   map[string]int `env:"MAP_STRING_INT"`
}

var (
	configOnce sync.Once
	cfg        *Config
)

func CreateConfig() *Config {
	if cfg != nil {
		return cfg
	}
	configOnce.Do(func() {
		cfg = &Config{}
		if err := env.Parse(cfg); err != nil {
			log.Error().Err(err).Msg("We have a problem with configuration!")
		}

	})

	return cfg
}
