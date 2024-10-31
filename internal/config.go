package internal

import (
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type swarmConfig struct {
	done         chan<- error
	seed         int64
	numMsgs      int
	engineConfig actor.EngineConfig
	msgTypes     []reflect.Type
}

type SwarmOption func(swarmConfig) swarmConfig

func WithSeed(seed int64) SwarmOption {
	return func(sc swarmConfig) swarmConfig {
		sc.seed = seed
		return sc
	}
}

func WithNumMsgs(numMsgs int) SwarmOption {
	return func(sc swarmConfig) swarmConfig {
		sc.numMsgs = numMsgs
		return sc
	}
}

func WithMessages(msgs ...any) SwarmOption {
	return func(sc swarmConfig) swarmConfig {
		for _, msg := range msgs {
			sc.msgTypes = append(sc.msgTypes, reflect.TypeOf(msg))
		}
		return sc
	}
}

func Run(initializer func(*actor.Engine) error, opts ...SwarmOption) error {
	done := make(chan error)
	config := swarmConfig{
		done:         done,
		seed:         time.Now().UnixNano(),
		numMsgs:      100,
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
