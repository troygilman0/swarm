package main

import (
	"log"
	"swarm/internal"

	"github.com/anthdm/hollywood/actor"
)

type MyType struct {
	Param1 string
	Param2 int
	Param3 bool
}

func main() {
	if err := internal.Run(initialize, internal.WithNumRounds(1)); err != nil {
		log.Println(err.Error())
	}
}

func initialize(engine *actor.Engine) error {
	for range 10 {
		engine.Spawn(testActorProducer(), "testActor")
	}
	return nil
}

type testActor struct{}

func testActorProducer() actor.Producer {
	return func() actor.Receiver {
		return &testActor{}
	}
}

func (a *testActor) Receive(act *actor.Context) {
	log.Printf("%+v", act.Message())
}
