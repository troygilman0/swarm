package swarm

import (
	"fmt"
	"log"
	"math/rand"
	"swarm/internal"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func Run(init actor.Producer, msgs []any, opts ...Option) error {
	opts = append(opts,
		withInitializer(init),
		withMessages(msgs),
	)

	config := options(opts).apply(Config{
		parallelRounds: 1,
	})

	var round uint64
	var roundsRunning uint64
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	results := make(chan result)

	for {
		if roundsRunning < config.parallelRounds {
			if config.numRounds > 0 && round >= config.numRounds && roundsRunning == 0 {
				break
			}
			done := make(chan error)
			roundConfig := options(opts).apply(Config{
				engineConfig: actor.NewEngineConfig(),
				SimulatorConfig: internal.SimulatorConfig{
					Done:     done,
					Seed:     random.Int63(),
					NumMsgs:  100,
					Interval: time.Millisecond,
				},
			})
			log.Printf("Starting round %d with seed %d\n", round, roundConfig.Seed)
			go runRound(roundConfig, done, results)
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

func runRound(config Config, done <-chan error, results chan<- result) {
	var err error
	start := time.Now()
	defer func() {
		results <- result{
			seed:     config.Seed,
			duration: time.Since(start),
			err:      err,
		}
	}()

	if config.SimulatorConfig.Interval == 0 {
		err = fmt.Errorf("interval cannot be 0")
		return
	}

	var engine *actor.Engine
	engine, err = actor.NewEngine(config.engineConfig)
	if err != nil {
		return
	}

	engine.Spawn(internal.NewSimulator(config.SimulatorConfig), "swarm-simulator")
	err = <-done
}
