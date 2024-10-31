package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type swarm struct {
	swarmConfig
	err      error
	rand     random
	msgCount int
	pids     []*actor.PID
}

func swarmProducer(config swarmConfig) actor.Producer {
	return func() actor.Receiver {
		return &swarm{
			swarmConfig: config,
		}
	}
}

func (s *swarm) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		s.err = nil
		s.msgCount = 0
		s.rand = newRandom(s.seed)
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
		s.done <- s.err

	case actor.ActorStartedEvent:
		if msg.PID == act.PID() {
			break
		}
		s.pids = append(s.pids, msg.PID)

	case actor.ActorRestartedEvent:
		if s.err != nil {
			break
		}
		s.err = fmt.Errorf("actor %s crashed at msg %d with seed %d", msg.PID.String(), s.msgCount, s.seed)
		act.Engine().Stop(act.PID())

	case sendMessagesMsg:
		if s.msgCount >= s.numMsgs {
			act.Engine().Stop(act.PID())
			break
		}
		pid := s.pids[s.rand.Intn(len(s.pids))]
		newMsgType := s.msgTypes[s.rand.Intn(len(s.msgTypes))]
		act.Send(pid, s.rand.Any(newMsgType))
		s.msgCount++
	}
}

type sendMessagesMsg struct{}

type TestMsg struct {
	Str string
}
