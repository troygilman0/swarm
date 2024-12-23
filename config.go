package swarm

import (
	"reflect"
	"swarm/internal/remoter"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type SwarmConfig struct {
	Initializer    actor.Producer
	msgTypes       []reflect.Type
	numMsgs        uint64
	numRounds      uint64
	parallelRounds uint64
	interval       time.Duration
	seed           int64
	adapter        remoter.Adapter
}

type Option func(SwarmConfig) SwarmConfig

func WithSeed(seed int64) Option {
	return func(c SwarmConfig) SwarmConfig {
		c.seed = seed
		return c
	}
}

func WithNumMessages(numMsgs uint64) Option {
	return func(c SwarmConfig) SwarmConfig {
		c.numMsgs = numMsgs
		return c
	}
}

func WithMessages(msgs []any) Option {
	return func(c SwarmConfig) SwarmConfig {
		for _, msg := range msgs {
			c.msgTypes = append(c.msgTypes, reflect.TypeOf(msg))
		}
		return c
	}
}

func WithNumRounds(numRounds uint64) Option {
	return func(c SwarmConfig) SwarmConfig {
		c.numRounds = numRounds
		return c
	}
}

func WithParellel(parallelRounds uint64) Option {
	return func(c SwarmConfig) SwarmConfig {
		c.parallelRounds = parallelRounds
		return c
	}
}

func WithInterval(interval time.Duration) Option {
	return func(c SwarmConfig) SwarmConfig {
		c.interval = interval
		return c
	}
}
