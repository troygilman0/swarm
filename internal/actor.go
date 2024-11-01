package internal

import (
	"fmt"
	"math/rand"
	"reflect"
	"swarm/internal/random"
	"sync"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type SwarmConfig struct {
	Done     chan<- error
	Seed     int64
	NumMsgs  uint64
	MsgTypes []reflect.Type
}

type swarm struct {
	SwarmConfig
	err      error
	rand     random.Random
	msgCount uint64
	pids     []*actor.PID
}

func NewSwarmProducer(config SwarmConfig) actor.Producer {
	return func() actor.Receiver {
		return &swarm{
			SwarmConfig: config,
		}
	}
}

func (s *swarm) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		s.err = nil
		s.msgCount = 0
		s.rand = random.NewRandom(rand.NewSource(s.Seed))
		s.pids = []*actor.PID{}

	case actor.Started:
		act.Engine().Subscribe(act.PID())
		act.SendRepeat(act.PID(), sendMessagesMsg{}, time.Millisecond)

	case actor.Stopped:
		wg := &sync.WaitGroup{}
		for _, pid := range s.pids {
			act.Engine().Stop(pid, wg)
		}
		wg.Wait()
		s.Done <- s.err

	case actor.ActorStartedEvent:
		if msg.PID == act.PID() {
			break
		}
		s.pids = append(s.pids, msg.PID)

	case actor.ActorRestartedEvent:
		if s.err != nil {
			break
		}
		s.err = fmt.Errorf("actor %s crashed at msg %d with seed %d", msg.PID.String(), s.msgCount, s.Seed)
		act.Engine().Stop(act.PID())

	case sendMessagesMsg:
		if s.msgCount >= s.NumMsgs {
			act.Engine().Stop(act.PID())
			break
		}
		pid := s.pids[s.rand.Intn(len(s.pids))]
		newMsgType := s.MsgTypes[s.rand.Intn(len(s.MsgTypes))]
		act.Send(pid, s.rand.Any(newMsgType))
		s.msgCount++
	}
}

type sendMessagesMsg struct{}
