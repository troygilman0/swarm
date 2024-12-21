package main

import (
	"log"
	"strings"
	"swarm"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func main() {
	if err := swarm.Run(
		newInitializer(),
		[]any{
			TestMsg{},
		},
		swarm.WithNumMsgs(10),
		swarm.WithParellel(10),
		swarm.WithInterval(time.Millisecond),
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
			act.Engine().Spawn(testActorProducer(), "testActor")
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
		if strings.Contains(msg.Str, "hel") {
			panic("56")
		}
	}
}
