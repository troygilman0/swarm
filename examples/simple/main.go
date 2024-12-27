package main

import (
	"log"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/swarm/simulator"
)

func main() {
	if err := simulator.Run(
		newInitializer(),
		simulator.WithMessages(
			TestMsg{},
		),
		simulator.WithNumMessages(100),
		simulator.WithParellel(10),
		simulator.WithInterval(time.Millisecond),
	); err != nil {
		log.Println(err)
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
			// act.Engine().Spawn(testActorProducer(), "testActor")
		}
	}
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
	switch msg := act.Message().(type) {
	case TestMsg:
		if strings.Contains(msg.Str, "he") {
			panic("")
		}
	}
}
