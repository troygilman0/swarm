package swarm

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/troygilman0/swarm/internal"

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
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	results := make(chan result)

	for {
		if roundsRunning < config.parallelRounds {
			if config.numRounds > 0 && round >= config.numRounds && roundsRunning == 0 {
				break
			}
			roundConfig := options(opts).apply(Config{
				engineConfig: actor.NewEngineConfig(),
				SimulatorConfig: internal.SimulatorConfig{
					Seed:     random.Int63(),
					NumMsgs:  100,
					Interval: time.Millisecond,
				},
			})
			if config.SimulatorConfig.Interval == 0 {
				panic("interval cannot be 0")
			}
			log.Printf("Starting round %d with seed %d\n", round, roundConfig.Seed)
			go runRound(roundConfig, results)
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

func runRound(config Config, results chan<- result) {
	var (
		start    = time.Now()
		done     = make(chan struct{})
		stopping bool
	)

	engine, err := actor.NewEngine(config.engineConfig)
	if err != nil {
		panic(err.Error())
	}

	engine.SpawnFunc(func(act *actor.Context) {
		switch msg := act.Message().(type) {
		case actor.Started:
			act.SpawnChild(internal.NewSimulatorProducer(config.SimulatorConfig), "swarm-simulator")

		case actor.Stopped:
			close(done)

		case internal.StopMsg:
			if !stopping {
				stopping = true
				engine.Stop(act.PID())
				results <- result{
					seed:     config.Seed,
					duration: time.Since(start),
					err:      msg.Err,
				}
			}
		}
	}, "swarm-adapter")

	cleanup := config.init()(engine)
	if cleanup != nil {
		defer cleanup()
	}

	<-done
}
