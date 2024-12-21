package swarm

import (
	"testing"

	"github.com/anthdm/hollywood/actor"
)

func testStarm(t *testing.T) {
	if err := Run(
		newInitializer(),
		[]any{testMsg{}},
		WithSeed(0),
		WithNumMsgs(100),
	); err != nil {
		t.Error(err)
	}
}

type initializer struct{}

func newInitializer() actor.Producer {
	return func() actor.Receiver {
		return &initializer{}
	}
}

func (i *initializer) Receive(act *actor.Context) {
	switch act.Message().(type) {
	case actor.Initialized:
		for range 10 {
			act.SpawnChild(testActorProducer(), "testActor")
		}
	}
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
