package swarm

import (
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func WithSeed(seed int64) Option {
	return func(c Config) Config {
		c.Seed = seed
		return c
	}
}

func WithNumMsgs(numMsgs uint64) Option {
	return func(c Config) Config {
		c.NumMsgs = numMsgs
		return c
	}
}

func withMessages(msgs []any) Option {
	return func(c Config) Config {
		for _, msg := range msgs {
			c.MsgTypes = append(c.MsgTypes, reflect.TypeOf(msg))
		}
		return c
	}
}

func withInitializer(init actor.Producer) Option {
	return func(c Config) Config {
		c.Initializer = init
		return c
	}
}

func WithNumRounds(numRounds uint64) Option {
	return func(c Config) Config {
		c.numRounds = numRounds
		return c
	}
}

func WithEngineConfig(config actor.EngineConfig) Option {
	return func(c Config) Config {
		c.engineConfig = config
		return c
	}
}

func WithParellel(parallelRounds uint64) Option {
	return func(c Config) Config {
		c.parallelRounds = parallelRounds
		return c
	}
}

func WithInterval(interval time.Duration) Option {
	return func(c Config) Config {
		c.Interval = interval
		return c
	}
}
