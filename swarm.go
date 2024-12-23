package swarm

import (
	"log"
	"math/rand"
	"strconv"
	"swarm/internal/remoter"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type swarm struct {
	SwarmConfig
	round        uint64
	activeRounds uint64
	random       *rand.Rand
	listenerPID  *actor.PID
	simulators   map[int64]*actor.PID
}

func newSwarm(config SwarmConfig) actor.Producer {
	return func() actor.Receiver {
		return &swarm{
			SwarmConfig: config,
		}
	}
}

func (s *swarm) Receive(act *actor.Context) {
	log.Printf("%T - %+v\n", act.Message(), act.Message())
	switch event := act.Message().(type) {
	case actor.Initialized:
		s.round = 0
		s.activeRounds = 0
		s.random = rand.New(rand.NewSource(time.Now().UnixNano()))
		s.listenerPID = nil
		s.simulators = make(map[int64]*actor.PID)

	case actor.Started:
		act.Send(act.PID(), startSimulationEvent{})

	case actor.Stopped:
		act.Send(s.listenerPID, swarmDoneEvent{})

	case registerListenerEvent:
		s.listenerPID = act.Sender()

	case startSimulationEvent:
		if s.activeRounds >= s.parallelRounds {
			return
		}
		seed := s.random.Int63()
		if _, ok := s.simulators[seed]; ok {
			// retry if seed is taken
			act.Send(act.PID(), startSimulationEvent{})
			return
		}
		address := strconv.FormatInt(seed, 10)
		engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter.NewLocalRemoter(s.adapter, address)))
		if err != nil {
			panic(err)
		}
		simulatorPID := engine.Spawn(newSimulator(simulatorConfig{
			swarmPID:    act.PID(),
			initializer: s.Initializer,
			seed:        seed,
			numMsgs:     s.numMsgs,
			interval:    s.interval,
			msgTypes:    s.msgTypes,
		}), "swarm-simulator")
		s.simulators[seed] = simulatorPID
		s.activeRounds++
		act.Send(act.PID(), startSimulationEvent{})

	case simulationDoneEvent:
		address := strconv.FormatInt(event.seed, 10)
		s.adapter.Stop(address).Wait()
		delete(s.simulators, event.seed)
		s.activeRounds--
		act.Send(act.PID(), startSimulationEvent{})

	case simulationErrorEvent:
		log.Println(event.err)
		act.Engine().Stop(act.Sender())
		act.Engine().Stop(act.PID())

	}
}
