package swarm

import (
	"fmt"
	"log"
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
		SimulatorConfig: internal.SimulatorConfig{
			Done:     done,
			Seed:     time.Now().UnixNano(),
			NumMsgs:  100,
			Interval: time.Millisecond,
		},
	})

	if config.SimulatorConfig.Interval == 0 {
		err = fmt.Errorf("interval cannot be 0")
		return
	}

	var engine *actor.Engine
	engine, err = actor.NewEngine(config.engineConfig)
	if err != nil {
		return
	}

	engine.Spawn(internal.NewSimulatorProducer(config.SimulatorConfig), "swarm-simulator")

	cleanup := config.init(engine)
	err = <-done
	if cleanup != nil {
		cleanup()
	}
}
