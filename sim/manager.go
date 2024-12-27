package sim

import (
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/troygilman0/swarm/internal/remoter"

	"github.com/anthdm/hollywood/actor"
)

type managerConfig struct {
	initializer    actor.Producer
	adapter        remoter.Adapter
	msgTypes       []reflect.Type
	numMsgs        uint64
	numRounds      uint64
	parallelRounds uint64
	interval       time.Duration
	seed           int64
}

const (
	StatusIdle = iota
	StatusRunning
	StatusStopping
)

type managerActor struct {
	managerConfig
	status       int
	round        uint64
	activeRounds uint64
	random       *rand.Rand
	listenerPID  *actor.PID
	simulators   map[int64]*actor.PID
}

func NewManager(initializer actor.Producer, adapter remoter.Adapter, opts ...Option) actor.Producer {
	config := managerConfig{
		initializer:    initializer,
		adapter:        adapter,
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
	// log.Printf("%s : %T - %+v\n", act.PID().String(), act.Message(), act.Message())
	switch msg := act.Message().(type) {
	case actor.Initialized:
		manager.status = StatusIdle
		manager.round = 0
		manager.activeRounds = 0
		manager.random = rand.New(rand.NewSource(time.Now().UnixNano()))
		manager.listenerPID = nil
		manager.simulators = make(map[int64]*actor.PID)

	case actor.Stopped:
		act.Send(manager.listenerPID, ManagerDoneEvent{})

	case Start:
		if manager.status != StatusIdle {
			return
		}
		manager.updateStatus(act, StatusRunning)
		act.Send(act.PID(), startSimulation{})

	case Stop:
		if manager.status != StatusRunning {
			return
		}
		manager.updateStatus(act, StatusStopping)
		for _, pid := range manager.simulators {
			act.Send(pid, stopSimulation{})
		}

	case RegisterListener:
		manager.listenerPID = act.Sender()

	case startSimulation:
		if manager.status != StatusRunning {
			return
		}
		if manager.activeRounds >= manager.parallelRounds {
			return
		}
		seed := manager.random.Int63()
		if _, ok := manager.simulators[seed]; ok {
			// retry if seed is taken
			act.Send(act.PID(), startSimulation{})
			return
		}
		address := strconv.FormatInt(seed, 10)
		engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter.NewRemoter(manager.adapter, address)))
		if err != nil {
			act.Send(manager.listenerPID, SimulationErrorEvent{
				Seed:  seed,
				Error: err,
			})
			return
		}
		simulatorPID := engine.Spawn(newSimulator(simulatorConfig{
			managerPID:  act.PID(),
			initializer: manager.initializer,
			seed:        seed,
			numMsgs:     manager.numMsgs,
			interval:    manager.interval,
			msgTypes:    manager.msgTypes,
		}), "swarm-simulator", actor.WithID(address))
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
		address := strconv.FormatInt(msg.Seed, 10)
		manager.adapter.Stop(address).Wait()
		delete(manager.simulators, msg.Seed)
		manager.activeRounds--
		act.Send(manager.listenerPID, msg)

		switch manager.status {
		case StatusRunning:
			act.Send(act.PID(), startSimulation{})
		case StatusStopping:
			if manager.activeRounds == 0 {
				manager.updateStatus(act, StatusIdle)
			}
		}

	case SimulationErrorEvent:
		act.Send(manager.listenerPID, msg)
		act.Engine().Stop(act.PID())

	}
}

func (manager *managerActor) updateStatus(act *actor.Context, status int) {
	manager.status = status
	act.Send(manager.listenerPID, ManagerStatusUpdateEvent{
		Status: manager.status,
	})
}
