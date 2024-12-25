package swarm

import (
	"reflect"
	"time"
)

type Option func(swarmConfig) swarmConfig

func WithSeed(seed int64) Option {
	return func(c swarmConfig) swarmConfig {
		c.seed = seed
		return c
	}
}

func WithNumMessages(numMsgs uint64) Option {
	return func(c swarmConfig) swarmConfig {
		c.numMsgs = numMsgs
		return c
	}
}

func WithMessages(msgs []any) Option {
	return func(c swarmConfig) swarmConfig {
		for _, msg := range msgs {
			c.msgTypes = append(c.msgTypes, reflect.TypeOf(msg))
		}
		return c
	}
}

func WithParellel(parallelRounds uint64) Option {
	return func(c swarmConfig) swarmConfig {
		if parallelRounds > 0 {
			c.parallelRounds = parallelRounds
		}
		return c
	}
}

func WithInterval(interval time.Duration) Option {
	return func(c swarmConfig) swarmConfig {
		if interval > 0 {
			c.interval = interval
		}
		return c
	}
}
