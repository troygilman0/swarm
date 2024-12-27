package main

import (
	"log"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/swarm/internal/remoter"
	"github.com/troygilman0/swarm/server"
	"github.com/troygilman0/swarm/sim"
)

func main() {
	adapter := remoter.NewLocalAdapter()
	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter.NewRemoter(adapter, "swarm")))
	if err != nil {
		panic(err)
	}

	pid := engine.Spawn(sim.NewManager(newInitializer(), adapter, sim.WithParellel(10)), "manager")

	engine.SpawnFunc(func(act *actor.Context) {
		log.Printf("%s : %T - %+v\n", act.PID().String(), act.Message(), act.Message())
		switch act.Message().(type) {
		case actor.Started:
			act.Send(pid, sim.RegisterListener{})
		}
	}, "listener")

	if err := http.ListenAndServe(":8080", server.NewSwarmHandler(engine, pid)); err != nil {
		panic(err)
	}
}

func newInitializer() actor.Producer {
	return func() actor.Receiver {
		return &initializer{}
	}
}

type initializer struct {
}

func (init *initializer) Receive(act *actor.Context) {

}
