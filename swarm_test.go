package swarm

import (
	"testing"

	"github.com/anthdm/hollywood/actor"
)

func TestStarm(t *testing.T) {
	if err := Run(
		initialize,
		WithSeed(0),
		WithNumMsgs(100),
		WithMessages(testMsg{}),
	); err != nil {
		t.Error(err)
	}
}

func initialize(engine *actor.Engine) func() {
	for range 10 {
		engine.Spawn(testActorProducer(), "testActor")
	}
	return nil
}

type testMsg struct {
	Str   string
	Int   int
	Uint  uint
	Float float64
	Bool  bool
	Slice []int
	Array [10]int
}

type testActor struct{}

func testActorProducer() actor.Producer {
	return func() actor.Receiver {
		return &testActor{}
	}
}

func (a *testActor) Receive(act *actor.Context) {
	switch act.Message().(type) {
	case actor.Initialized, actor.Started, actor.Stopped:
	default:
		// log.Printf("%T - %+v", act.Message(), act.Message())
	}
}
