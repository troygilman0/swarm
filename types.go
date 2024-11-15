package swarm

import (
	"time"

	"github.com/troygilman0/swarm/internal"

	"github.com/anthdm/hollywood/actor"
)

type Initializer func() func(*actor.Engine) func()

type result struct {
	seed     int64
	duration time.Duration
	err      error
}

type Config struct {
	init           Initializer
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
