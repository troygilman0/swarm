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
	baseSeed := time.Now().UnixNano()
	results := make(chan result)

	for {
		if roundsRunning < config.parallelRounds {
			if config.numRounds > 0 && round >= config.numRounds && roundsRunning == 0 {
				break
			}
			seed := baseSeed + int64(round)
			log.Printf("Starting round %d with seed %d\n", round, seed)
			go runRound(opts, seed, results)
			round++
			roundsRunning++
		} else {
			result := <-results
			if result.err != nil {
				return fmt.Errorf("seed %d resulted in error: %s", result.seed, result.err.Error())
			}
			log.Printf("Round with seed %d finished in %f seconds\n", result.seed, result.duration.Seconds())
			roundsRunning--
		}
	}

	return nil
}

func runRound(opts options, seed int64, results chan<- result) {
	var err error
	start := time.Now()
	defer func() {
		results <- result{
			seed:     seed,
			duration: time.Since(start),
			err:      err,
		}
	}()

	done := make(chan error)
	config := opts.apply(Config{
		engineConfig: actor.NewEngineConfig(),
		SimulatorConfig: internal.SimulatorConfig{
			Done:     done,
			Seed:     seed,
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
