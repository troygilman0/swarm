package sim

import (
	"github.com/anthdm/hollywood/actor"
)

func Run(initializer actor.Producer, opts ...Option) error {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		return err
	}

	managerPID := engine.Spawn(NewManager(initializer, opts...), "swarm-manager")

	done := make(chan struct{})
	engine.SpawnFunc(func(act *actor.Context) {
		switch act.Message().(type) {
		case actor.Started:
			act.Send(managerPID, RegisterListener{})
		case ManagerDoneEvent:
			close(done)
		}
	}, "swarm-listener")

	<-done
	return nil
}
