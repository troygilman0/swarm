package main

import (
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

	pid := engine.Spawn(sim.NewManager(newInitializer(), adapter), "manager")

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
