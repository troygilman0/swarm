package internal

import (
	"time"

	"github.com/anthdm/hollywood/actor"
)

type swarmConfig struct {
	done         chan<- error
	seed         int64
	numRounds    int
	engineConfig actor.EngineConfig
}

type SwarmOption func(swarmConfig) swarmConfig

func WithSeed(seed int64) SwarmOption {
	return func(sc swarmConfig) swarmConfig {
		sc.seed = seed
		return sc
	}
}

func WithNumRounds(numRounds int) SwarmOption {
	return func(sc swarmConfig) swarmConfig {
		sc.numRounds = numRounds
		return sc
	}
}

func WithMessage(msg any) SwarmOption {
	return func(sc swarmConfig) swarmConfig {
		return sc
	}
}

func Run(initializer func(*actor.Engine) error, opts ...SwarmOption) error {
	done := make(chan error)
	config := swarmConfig{
		done:         done,
		seed:         time.Now().UnixNano(),
		numRounds:    100,
		engineConfig: actor.NewEngineConfig(),
	}

	for _, opt := range opts {
		config = opt(config)
	}

	engine, err := actor.NewEngine(config.engineConfig)
	if err != nil {
		return err
	}

	engine.Spawn(swarmProducer(config), "swarm")

	if err := initializer(engine); err != nil {
		return err
	}

	err = <-done
	return err
}
