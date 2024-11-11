package internal

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/troygilman0/swarm/internal/random"

	"github.com/anthdm/hollywood/actor"
)

type SimulatorConfig struct {
	Seed     int64
	NumMsgs  uint64
	Interval time.Duration
	MsgTypes []reflect.Type
}

type simulator struct {
	SimulatorConfig
	err      error
	rand     random.Random
	msgCount uint64
	pids     []*actor.PID
	pidsSet  map[*actor.PID]struct{}
	repeater actor.SendRepeater
}

func NewSimulatorProducer(config SimulatorConfig) actor.Producer {
	return func() actor.Receiver {
		return &simulator{
			SimulatorConfig: config,
		}
	}
}

func (s *simulator) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
		s.err = nil
		s.msgCount = 0
		s.rand = random.NewRandom(rand.NewSource(s.Seed))
		s.pids = []*actor.PID{}
		s.pidsSet = make(map[*actor.PID]struct{})

	case actor.Started:
		act.Engine().Subscribe(act.PID())
		s.repeater = act.SendRepeat(act.PID(), sendRandomMsg{}, s.Interval)

	case actor.Stopped:
		act.Engine().Unsubscribe(act.PID())
		s.repeater.Stop()
		wg := &sync.WaitGroup{}
		for _, pid := range s.pids {
			act.Engine().Stop(pid, wg)
		}
		wg.Wait()

	case actor.ActorStartedEvent:
		if msg.PID == act.PID() || msg.PID == act.Parent() {
			break
		}
		s.pids = append(s.pids, msg.PID)
		s.pidsSet[msg.PID] = struct{}{}

	case actor.ActorRestartedEvent:
		if _, ok := s.pidsSet[msg.PID]; !ok {
			break
		}
		act.Send(act.Parent(), StopMsg{
			Err: fmt.Errorf("actor %s crashed at msg %d", msg.PID.String(), s.msgCount),
		})

	case sendRandomMsg:
		if s.msgCount >= s.NumMsgs {
			act.Send(act.Parent(), StopMsg{})
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

type StopMsg struct {
	Err error
}

type sendRandomMsg struct{}
