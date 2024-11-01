package swarm

import (
	"log"
	"reflect"
	"swarm/internal"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func Run(init func(*actor.Engine), opts ...Option) error {
	config := options(opts).apply(Config{})

	var round uint64
	for {
		if config.numRounds > 0 && round >= config.numRounds {
			break
		}
		start := time.Now()
		if err := runRound(init, opts); err != nil {
			return err
		}
		log.Printf("Ran round %d in %f seconds", round, time.Since(start).Seconds())
		round++
	}

	return nil
}

func runRound(init func(*actor.Engine), opts options) error {
	done := make(chan error)
	config := Config{
		engineConfig: actor.NewEngineConfig(),
		SwarmConfig: internal.SwarmConfig{
			Done:    done,
			Seed:    time.Now().UnixNano(),
			NumMsgs: 100,
		},
	}

	config = options(opts).apply(config)

	engine, err := actor.NewEngine(config.engineConfig)
	if err != nil {
		return err
	}

	engine.Spawn(internal.NewSwarmProducer(config.SwarmConfig), "swarm")

	init(engine)

	return <-done
}

type Config struct {
	engineConfig actor.EngineConfig
	numRounds    uint64
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

func WithMessages(msgs ...any) Option {
	return func(c Config) Config {
		for _, msg := range msgs {
			c.MsgTypes = append(c.MsgTypes, reflect.TypeOf(msg))
		}
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
