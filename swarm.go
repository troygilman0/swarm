package swarm

import (
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/troygilman0/swarm/internal/remoter"

	"github.com/anthdm/hollywood/actor"
)

type swarmConfig struct {
	initializer    actor.Producer
	adapter        remoter.Adapter
	msgTypes       []reflect.Type
	numMsgs        uint64
	numRounds      uint64
	parallelRounds uint64
	interval       time.Duration
	seed           int64
}

type swarm struct {
	swarmConfig
	round        uint64
	activeRounds uint64
	random       *rand.Rand
	listenerPID  *actor.PID
	simulators   map[int64]*actor.PID
}

func NewSwarm(initializer actor.Producer, adapter remoter.Adapter, opts ...Option) actor.Producer {
	config := swarmConfig{
		initializer:    initializer,
		adapter:        adapter,
		parallelRounds: 1,
		interval:       time.Millisecond,
	}

	for _, opt := range opts {
		config = opt(config)
	}

	return func() actor.Receiver {
		return &swarm{
			swarmConfig: config,
		}
	}
}

func (s *swarm) Receive(act *actor.Context) {
	log.Printf("%s : %T - %+v\n", act.PID().String(), act.Message(), act.Message())
	switch msg := act.Message().(type) {
	case actor.Initialized:
		s.round = 0
		s.activeRounds = 0
		s.random = rand.New(rand.NewSource(time.Now().UnixNano()))
		s.listenerPID = nil
		s.simulators = make(map[int64]*actor.PID)

	case actor.Started:
		act.Send(act.PID(), StartSimulation{})

	case actor.Stopped:
		act.Send(s.listenerPID, SwarmDoneEvent{})

	case RegisterListener:
		s.listenerPID = act.Sender()

	case StartSimulation:
		if s.activeRounds >= s.parallelRounds {
			return
		}
		seed := s.random.Int63()
		if _, ok := s.simulators[seed]; ok {
			// retry if seed is taken
			act.Send(act.PID(), StartSimulation{})
			return
		}
		address := strconv.FormatInt(seed, 10)
		engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter.NewRemoter(s.adapter, address)))
		if err != nil {
			act.Send(s.listenerPID, SimulationErrorEvent{
				Seed:  seed,
				Error: err,
			})
			return
		}
		simulatorPID := engine.Spawn(newSimulator(simulatorConfig{
			swarmPID:    act.PID(),
			initializer: s.initializer,
			seed:        seed,
			numMsgs:     s.numMsgs,
			interval:    s.interval,
			msgTypes:    s.msgTypes,
		}), "swarm-simulator", actor.WithID(address))
		s.simulators[seed] = simulatorPID
		s.activeRounds++
		act.Send(act.PID(), StartSimulation{})

	case SimulationDoneEvent:
		if _, ok := s.simulators[msg.Seed]; !ok {
			return
		}
		address := strconv.FormatInt(msg.Seed, 10)
		s.adapter.Stop(address).Wait()
		delete(s.simulators, msg.Seed)
		s.activeRounds--
		act.Send(act.PID(), StartSimulation{})

	case SimulationErrorEvent:
		act.Send(s.listenerPID, msg)
		act.Engine().Stop(act.PID())

	}
}
