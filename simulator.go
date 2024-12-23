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
	rand     random.Random
	msgCount uint64
	pids     []*actor.PID
	repeater actor.SendRepeater
}

func newSimulator(config simulatorConfig) actor.Producer {
	return func() actor.Receiver {
		return &simulator{
			simulatorConfig: config,
		}
	}
}

func (s *simulator) Receive(act *actor.Context) {
	// log.Printf("%T - %+v\n", act.Message(), act.Message())
	switch msg := act.Message().(type) {
	case actor.Initialized:
		s.msgCount = 0
		s.rand = random.NewRandom(rand.NewSource(s.seed))
		s.pids = []*actor.PID{}

	case actor.Started:
		act.Engine().Subscribe(act.PID())
		act.SpawnChild(s.initializer, "swarm-initializer")
		s.repeater = act.SendRepeat(act.PID(), sendMessagesMsg{}, s.interval)

	case actor.Stopped:
		act.Engine().Unsubscribe(act.PID())
		s.repeater.Stop()
		act.Send(s.swarmPID, simulationDoneEvent{
			seed: s.seed,
		})

	case actor.ActorStartedEvent:
		if msg.PID == act.PID() {
			break
		}
		s.pids = append(s.pids, msg.PID)

	case actor.ActorRestartedEvent:
		if msg.PID == act.PID() {
			break
		}
		act.Send(s.swarmPID, simulationErrorEvent{
			seed: s.seed,
			err:  fmt.Errorf("actor %s crashed at msg %d", msg.PID.String(), s.msgCount),
		})

	case sendMessagesMsg:
		if s.msgCount >= s.numMsgs {
			act.Engine().Stop(act.PID())
			break
		}
		if len(s.pids) == 0 {
			break
		}
		pid := s.pids[s.rand.Intn(len(s.pids))]
		newMsgType := s.msgTypes[s.rand.Intn(len(s.msgTypes))]
		act.Send(pid, s.rand.Any(newMsgType))
		s.msgCount++
	}
}

type sendMessagesMsg struct{}
