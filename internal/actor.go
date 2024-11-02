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
	Interval time.Duration
	MsgTypes []reflect.Type
}

type swarm struct {
	SwarmConfig
	stopping bool
	err      error
	rand     random.Random
	msgCount uint64
	pids     []*actor.PID
	repeater actor.SendRepeater
}

func NewSwarmProducer(config SwarmConfig) actor.Producer {
	return func() actor.Receiver {
		return &swarm{
			SwarmConfig: config,
		}
	}
}

func (s *swarm) Receive(act *actor.Context) {
	if _, ok := act.Message().(actor.Stopped); !ok && s.stopping {
		return
	}

	switch msg := act.Message().(type) {
	case actor.Initialized:
		s.err = nil
		s.msgCount = 0
		s.rand = random.NewRandom(rand.NewSource(s.Seed))
		s.pids = []*actor.PID{}

	case actor.Started:
		act.Engine().Subscribe(act.PID())
		s.repeater = act.SendRepeat(act.PID(), sendMessagesMsg{}, s.Interval)

	case actor.Stopped:
		act.Engine().Unsubscribe(act.PID())
		s.repeater.Stop()
		wg := &sync.WaitGroup{}
		for _, pid := range s.pids {
			act.Engine().Stop(pid, wg)
		}
		wg.Wait()
		s.Done <- s.err
		close(s.Done)

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
		s.stopping = true
		act.Engine().Stop(act.PID())

	case sendMessagesMsg:
		if s.msgCount >= s.NumMsgs {
			s.stopping = true
			act.Engine().Stop(act.PID())
			break
		}
		if len(s.pids) == 0 {
			break
		}
		pid := s.pids[s.rand.Intn(len(s.pids))]
		newMsgType := s.MsgTypes[s.rand.Intn(len(s.MsgTypes))]
		act.Send(pid, s.rand.Any(newMsgType))
		s.msgCount++
	}
}

type sendMessagesMsg struct{}
