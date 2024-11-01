package main

import (
	"log"
	"swarm"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func main() {
	if err := swarm.Run(
		initialize,
		// swarm.WithSeed(0),
		swarm.WithNumMsgs(1000),
		swarm.WithMessages(&TestMsg{}),
		swarm.WithParellel(10),
		swarm.WithInterval(time.Nanosecond),
	); err != nil {
		log.Println(err)
	}
}

func initialize(engine *actor.Engine) func() {
	for range 10 {
		engine.Spawn(testActorProducer(), "testActor")
	}
	return nil
}

type TestMsg struct {
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
