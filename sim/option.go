package sim

import (
	"reflect"
	"time"
)

type Option func(managerConfig) managerConfig

func WithSeed(seed int64) Option {
	return func(c managerConfig) managerConfig {
		c.seed = seed
		return c
	}
}

func WithNumMessages(numMsgs uint64) Option {
	return func(c managerConfig) managerConfig {
		c.numMsgs = numMsgs
		return c
	}
}

func WithMessages(msgs ...any) Option {
	return func(c managerConfig) managerConfig {
		for _, msg := range msgs {
			c.msgTypes = append(c.msgTypes, reflect.TypeOf(msg))
		}
		return c
	}
}

func WithParellel(parallelRounds uint64) Option {
	return func(c managerConfig) managerConfig {
		if parallelRounds > 0 {
			c.parallelRounds = parallelRounds
		}
		return c
	}
}

func WithInterval(interval time.Duration) Option {
	return func(c managerConfig) managerConfig {
		if interval > 0 {
			c.interval = interval
		}
		return c
	}
}
