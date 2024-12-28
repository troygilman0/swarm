package sim

import (
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type managerConfig struct {
	initializer    actor.Producer
	msgTypes       []reflect.Type
	numMsgs        uint64
	numRounds      uint64
	parallelRounds uint64
	interval       time.Duration
	seed           int64
}

type managerActor struct {
	managerConfig
	round        uint64
	activeRounds uint64
	random       *rand.Rand
	listenerPID  *actor.PID
	simulators   map[int64]*actor.PID
}

func NewManager(initializer actor.Producer, opts ...Option) actor.Producer {
	config := managerConfig{
		initializer:    initializer,
		parallelRounds: 1,
		interval:       time.Millisecond,
	}

	for _, opt := range opts {
		config = opt(config)
	}

	return func() actor.Receiver {
		return &managerActor{
			managerConfig: config,
		}
	}
}

func (manager *managerActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		manager.round = 0
		manager.activeRounds = 0
		manager.random = rand.New(rand.NewSource(time.Now().UnixNano()))
		manager.listenerPID = nil
		manager.simulators = make(map[int64]*actor.PID)

	case actor.Started:
		act.Send(act.PID(), startSimulation{})

	case RegisterListener:
		manager.listenerPID = act.Sender()

	case startSimulation:
		if manager.activeRounds >= manager.parallelRounds {
			return
		}
		seed := manager.random.Int63()
		if _, ok := manager.simulators[seed]; ok {
			// retry if seed is taken
			act.Send(act.PID(), startSimulation{})
			return
		}
		simulatorPID := act.SpawnChild(newSimulator(simulatorConfig{
			initializer: manager.initializer,
			seed:        seed,
			numMsgs:     manager.numMsgs,
			interval:    manager.interval,
			msgTypes:    manager.msgTypes,
		}), "swarm-simulator", actor.WithID(strconv.FormatInt(seed, 10)))
		manager.simulators[seed] = simulatorPID
		manager.activeRounds++
		act.Send(act.PID(), startSimulation{})

	case SimulationStartedEvent:
		if _, ok := manager.simulators[msg.Seed]; !ok {
			return
		}
		act.Send(manager.listenerPID, msg)

	case SimulationDoneEvent:
		if _, ok := manager.simulators[msg.Seed]; !ok {
			return
		}
		delete(manager.simulators, msg.Seed)
		manager.activeRounds--
		act.Send(manager.listenerPID, msg)
		act.Send(act.PID(), startSimulation{})

	case SimulationErrorEvent:
		act.Send(manager.listenerPID, msg)
		act.Engine().Stop(act.PID())

	}
}
