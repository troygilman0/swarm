package sim

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/troygilman0/swarm/internal/random"

	"github.com/anthdm/hollywood/actor"
)

type simulatorConfig struct {
	initializer actor.Producer
	seed        int64
	numMsgs     uint64
	interval    time.Duration
	msgTypes    []reflect.Type
}

type simulatorActor struct {
	simulatorConfig
	rand           random.Random
	msgCount       uint64
	engine         *actor.Engine
	pids           []*actor.PID
	bridgePID      *actor.PID
	initializerPID *actor.PID
}

func newSimulator(config simulatorConfig) actor.Producer {
	return func() actor.Receiver {
		return &simulatorActor{
			simulatorConfig: config,
		}
	}
}

func (simulator *simulatorActor) Receive(act *actor.Context) {
	// log.Printf("%s : %T - %+v\n", act.PID().String(), act.Message(), act.Message())
	switch msg := act.Message().(type) {
	case actor.Initialized:
		simulator.msgCount = 0
		simulator.rand = random.NewRandom(rand.NewSource(simulator.seed))
		simulator.engine, _ = actor.NewEngine(actor.NewEngineConfig())
		simulator.pids = []*actor.PID{}
		simulator.bridgePID = nil
		simulator.initializerPID = nil

	case actor.Started:
		simulator.bridgePID = simulator.engine.Spawn(newBridgeActor(act.Engine(), act.PID()), "bridge")
		simulator.initializerPID = simulator.engine.Spawn(simulator.initializer, "swarm-initializer")
		act.Send(act.Parent(), SimulationStartedEvent{Seed: simulator.seed})
		act.Send(act.PID(), sendMessagesMsg{})

	case actor.Stopped:
		simulator.engine.Stop(simulator.initializerPID).Wait()
		simulator.engine.Stop(simulator.bridgePID).Wait()
		act.Send(act.Parent(), SimulationDoneEvent{
			Seed: simulator.seed,
		})

	case stopSimulation:
		act.Engine().Stop(act.PID())

	case simulationEvent:
		switch event := msg.event.(type) {
		case actor.ActorStartedEvent:
			simulator.pids = append(simulator.pids, event.PID)

		case actor.ActorRestartedEvent:
			act.Send(act.Parent(), SimulationErrorEvent{
				Seed:  simulator.seed,
				Error: fmt.Errorf("actor %s crashed at msg %d", event.PID.String(), simulator.msgCount),
			})
			act.Engine().Stop(act.PID())
		}

	case sendMessagesMsg:
		if simulator.msgCount >= simulator.numMsgs {
			act.Engine().Stop(act.PID())
			return
		}
		sendWithDelay(act, act.PID(), sendMessagesMsg{}, simulator.interval)
		if len(simulator.pids) == 0 {
			return
		}
		pid := simulator.pids[simulator.rand.Intn(len(simulator.pids))]
		newMsgType := simulator.msgTypes[simulator.rand.Intn(len(simulator.msgTypes))]
		simulator.engine.Send(pid, simulator.rand.Any(newMsgType))
		simulator.msgCount++
	}
}

type sendMessagesMsg struct{}
