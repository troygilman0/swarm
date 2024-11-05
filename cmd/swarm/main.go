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
		initialize,
		[]any{
			TestMsg{},
		},
		swarm.WithNumMsgs(1000),
		swarm.WithParellel(10),
		swarm.WithInterval(time.Millisecond),
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
	switch msg := act.Message().(type) {
	case TestMsg:
		if strings.Contains(msg.Str, "hel") {
			panic("56")
		}
	}
}
