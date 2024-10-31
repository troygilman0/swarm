package main

import (
	"log"
	"swarm/internal"

	"github.com/anthdm/hollywood/actor"
)

func main() {
	if err := internal.Run(
		initialize,
		internal.WithNumMsgs(1000),
		internal.WithMessages(testMsg{}),
	); err != nil {
		log.Println(err.Error())
	}
}

func initialize(engine *actor.Engine) error {
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
		log.Printf("%T - %+v", act.Message(), act.Message())
	}
}
