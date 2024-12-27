package sim

import (
	"github.com/troygilman0/swarm/internal/remoter"

	"github.com/anthdm/hollywood/actor"
)

func Run(initializer actor.Producer, opts ...Option) error {
	adapter := remoter.NewLocalAdapter()

	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter.NewRemoter(adapter, "swarm")))
	if err != nil {
		return err
	}

	managerPID := engine.Spawn(NewManager(initializer, adapter, opts...), "swarm-manager")

	done := make(chan struct{})
	engine.SpawnFunc(func(act *actor.Context) {
		switch act.Message().(type) {
		case actor.Started:
			act.Send(managerPID, RegisterListener{})
		case SwarmDoneEvent:
			close(done)
		}
	}, "swarm-listener")

	<-done
	return nil
}
