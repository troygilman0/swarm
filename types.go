package swarm

import (
	"swarm/internal"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type Initializer func(*actor.Engine) func()

type result struct {
	seed     int64
	duration time.Duration
	err      error
}

type Config struct {
	engineConfig   actor.EngineConfig
	numRounds      uint64
	parallelRounds uint64
	internal.SimulatorConfig
}

type Option func(Config) Config

type options []Option

func (opts options) apply(config Config) Config {
	for _, opt := range opts {
		config = opt(config)
	}
	return config
}
