package swarm

import (
	"fmt"
	"math/rand"
	"reflect"
	"swarm/internal/random"
	"time"

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

type simulator struct {
	simulatorConfig
	rand           random.Random
	msgCount       uint64
	pids           []*actor.PID
	initializerPID *actor.PID
}

func newSimulator(config simulatorConfig) actor.Producer {
	return func() actor.Receiver {
		return &simulator{
			simulatorConfig: config,
		}
	}
}

func (s *simulator) Receive(act *actor.Context) {
	// log.Printf("%s : %T - %+v\n", act.PID().String(), act.Message(), act.Message())
	switch msg := act.Message().(type) {
	case actor.Initialized:
		s.msgCount = 0
		s.rand = random.NewRandom(rand.NewSource(s.seed))
		s.pids = []*actor.PID{}

	case actor.Started:
		act.Engine().Subscribe(act.PID())
		s.initializerPID = act.Engine().Spawn(s.initializer, "swarm-initializer")
		act.Send(act.PID(), sendMessagesMsg{})

	case actor.Stopped:
		act.Engine().Unsubscribe(act.PID())
		act.Engine().Stop(s.initializerPID).Wait()
		act.Send(s.swarmPID, simulationDoneEvent{
			seed: s.seed,
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
		s.pids = append(s.pids, msg.PID)

	case actor.ActorRestartedEvent:
		act.Send(s.swarmPID, simulationErrorEvent{
			seed: s.seed,
			err:  fmt.Errorf("actor %s crashed at msg %d", msg.PID.String(), s.msgCount),
		})
		act.Engine().Stop(act.PID())

	case sendMessagesMsg:
		if s.msgCount >= s.numMsgs {
			// act.Engine().Stop(s.initializerPID)
			// log.Println("SENDING TO", s.swarmPID)
			// act.Send(s.swarmPID, simulationDoneEvent{
			// 	seed: s.seed,
			// })
			act.Engine().Stop(act.PID())
			return
		}
		sendWithDelay(act, act.PID(), sendMessagesMsg{}, s.interval)
		if len(s.pids) == 0 {
			return
		}
		pid := s.pids[s.rand.Intn(len(s.pids))]
		newMsgType := s.msgTypes[s.rand.Intn(len(s.msgTypes))]
		act.Send(pid, s.rand.Any(newMsgType))
		s.msgCount++
	}
}

type sendMessagesMsg struct{}
