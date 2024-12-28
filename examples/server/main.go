package main

import (
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/swarm/server"
	"github.com/troygilman0/swarm/sim"
)

func main() {
	if err := http.ListenAndServe(":8080", server.NewSwarmHandler(newInitializer(), sim.WithParellel(10))); err != nil {
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
