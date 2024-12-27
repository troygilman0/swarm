package simulator

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/troygilman0/swarm/internal/random"

	"github.com/anthdm/hollywood/actor"
)

type simulatorConfig struct {
	swarmPID    *actor.PID
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
	pids           []*actor.PID
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
		simulator.pids = []*actor.PID{}

	case actor.Started:
		act.Engine().Subscribe(act.PID())
		simulator.initializerPID = act.Engine().Spawn(simulator.initializer, "swarm-initializer")
		act.Send(act.PID(), sendMessagesMsg{})

	case actor.Stopped:
		act.Engine().Unsubscribe(act.PID())
		act.Engine().Stop(simulator.initializerPID).Wait()
		act.Send(simulator.swarmPID, SimulationDoneEvent{
			Seed: simulator.seed,
		})

	case actor.DeadLetterEvent:
		panic(msg)

	case actor.ActorStartedEvent:
		if msg.PID == act.PID() {
			return
		}
		if act.Child(msg.PID.GetID()) != nil {
			return
		}
		simulator.pids = append(simulator.pids, msg.PID)

	case actor.ActorRestartedEvent:
		act.Send(simulator.swarmPID, SimulationErrorEvent{
			Seed:  simulator.seed,
			Error: fmt.Errorf("actor %s crashed at msg %d", msg.PID.String(), simulator.msgCount),
		})
		act.Engine().Stop(act.PID())

	case sendMessagesMsg:
		if simulator.msgCount >= simulator.numMsgs {
			// act.Engine().Stop(s.initializerPID)
			// log.Println("SENDING TO", s.swarmPID)
			// act.Send(s.swarmPID, simulationDoneEvent{
			// 	seed: s.seed,
			// })
			act.Engine().Stop(act.PID())
			return
		}
		sendWithDelay(act, act.PID(), sendMessagesMsg{}, simulator.interval)
		if len(simulator.pids) == 0 {
			return
		}
		pid := simulator.pids[simulator.rand.Intn(len(simulator.pids))]
		newMsgType := simulator.msgTypes[simulator.rand.Intn(len(simulator.msgTypes))]
		act.Send(pid, simulator.rand.Any(newMsgType))
		simulator.msgCount++
	}
}

type sendMessagesMsg struct{}
