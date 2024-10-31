package internal

import (
	"fmt"
	"sync"
	"time"

	"math/rand"

	"github.com/anthdm/hollywood/actor"
)

type swarm struct {
	swarmConfig
	err    error
	random *rand.Rand
	round  int
	pids   []*actor.PID
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
		s.round = 0
		s.random = rand.New(rand.NewSource(s.seed))
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
		s.err = fmt.Errorf("actor %s crashed at round %d with seed %d", msg.PID.String(), s.round, s.seed)
		act.Engine().Stop(act.PID())

	case sendMessagesMsg:
		if s.round >= s.numRounds {
			act.Engine().Stop(act.PID())
			break
		}
		for _, pid := range s.pids {
			act.Send(pid, randomize(TestMsg{}, s.random))
		}
		s.round++
	}
}

type sendMessagesMsg struct{}

type TestMsg struct {
	Str string
}
