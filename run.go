package swarm

import (
	"swarm/internal/remoter"

	"github.com/anthdm/hollywood/actor"
)

func Run(initializer actor.Producer, opts ...Option) error {
	adapter := remoter.NewLocalAdapter()

	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter.NewRemoter(adapter, "swarm")))
	if err != nil {
		return err
	}

	swarmPID := engine.Spawn(NewSwarm(initializer, adapter, opts...), "swarm")

	done := make(chan struct{})
	engine.SpawnFunc(func(act *actor.Context) {
		switch act.Message().(type) {
		case actor.Started:
			act.Send(swarmPID, RegisterListener{})
		case SwarmDoneEvent:
			close(done)
		}
	}, "swarm-listener")

	<-done
	return nil
}
