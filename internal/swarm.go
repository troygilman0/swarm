package internal

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"math/rand"

	"github.com/anthdm/hollywood/actor"
)

type Swarm struct {
	producers []producer
}

func NewSwarm(opts ...SwarmOption) Swarm {
	swarm := Swarm{}
	for _, opt := range opts {
		opt(swarm)
	}
	return swarm
}

type SwarmOption func(swarm Swarm)

func WithMessage(msg any) SwarmOption {
	return func(swarm Swarm) {
		v := reflect.ValueOf(msg)
		swarm.producers = append(swarm.producers, newProducer(v))
	}
}

func newProducer(v reflect.Value) producer {
	return func() any {
		if v.Kind() != reflect.Ptr || v.IsNil() {
			return nil
		}

		v = v.Elem()
		if v.Kind() != reflect.Struct {
			return nil
		}

		fields := []reflect.Value{}
		for i := range v.NumField() {
			field := v.Field(i)
			if !field.IsValid() {
				return nil
			}
			if !field.CanSet() {
				return nil
			}

			fields = append(fields, field)
		}
		return nil
	}
}

func populateValue(v reflect.Value) {

}

type producer func() any

func Run(config actor.EngineConfig, initializer func(*actor.Engine) error, seed *int64) error {
	if seed == nil {
		s := time.Now().UnixNano()
		seed = &s
	}

	engine, err := actor.NewEngine(config)
	if err != nil {
		return err
	}

	done := make(chan error)
	engine.Spawn(swarmProducer(done, *seed, 5), "swarm")

	if err := initializer(engine); err != nil {
		return err
	}

	err = <-done
	return err
}

type swarm struct {
	done      chan<- error
	err       error
	seed      int64
	random    *rand.Rand
	numRounds int
	round     int
	pids      []*actor.PID
}

func swarmProducer(done chan<- error, seed int64, numRounds int) actor.Producer {
	return func() actor.Receiver {
		return &swarm{
			done:      done,
			seed:      seed,
			numRounds: numRounds,
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
			act.Send(pid, s.random.Int())
		}
		s.round++
	}
}

type sendMessagesMsg struct{}
