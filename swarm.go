package swarm

import (
	"fmt"
	"log"
	"reflect"
	"swarm/internal"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func Run(init Initializer, msgs []any, opts ...Option) error {
	opts = append(opts,
		withInitializer(init),
		withMessages(msgs),
	)

	config := options(opts).apply(Config{
		parallelRounds: 1,
	})

	var round uint64
	var roundsRunning uint64
	results := make(chan result)

	for {
		if roundsRunning < config.parallelRounds {
			if config.numRounds > 0 && round >= config.numRounds && roundsRunning == 0 {
				break
			}
			log.Printf("Starting round %d\n", round)
			go runRound(opts, round, results)
			round++
			roundsRunning++
		} else {
			result := <-results
			if result.err != nil {
				return result.err
			}
			log.Printf("Run round %d in %f seconds\n", result.round, result.duration.Seconds())
			roundsRunning--
		}
	}

	return nil
}

func runRound(opts options, round uint64, results chan<- result) {
	var err error
	start := time.Now()
	defer func() {
		results <- result{
			round:    round,
			duration: time.Since(start),
			err:      err,
		}
	}()

	done := make(chan error)
	config := opts.apply(Config{
		engineConfig: actor.NewEngineConfig(),
		SwarmConfig: internal.SwarmConfig{
			Done:     done,
			Seed:     time.Now().UnixNano(),
			NumMsgs:  100,
			Interval: time.Millisecond,
		},
	})

	if config.SwarmConfig.Interval == 0 {
		err = fmt.Errorf("interval cannot be 0")
		return
	}

	var engine *actor.Engine
	engine, err = actor.NewEngine(config.engineConfig)
	if err != nil {
		return
	}

	engine.Spawn(internal.NewSwarmProducer(config.SwarmConfig), "swarm")

	cleanup := config.init(engine)
	err = <-done
	if cleanup != nil {
		cleanup()
	}
}

type Initializer func(*actor.Engine) func()

type result struct {
	round    uint64
	duration time.Duration
	err      error
}

type Config struct {
	init           Initializer
	engineConfig   actor.EngineConfig
	numRounds      uint64
	parallelRounds uint64
	internal.SwarmConfig
}

type Option func(Config) Config

type options []Option

func (opts options) apply(config Config) Config {
	for _, opt := range opts {
		config = opt(config)
	}
	return config
}

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

func withInitializer(init Initializer) Option {
	return func(c Config) Config {
		c.init = init
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
