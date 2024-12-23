package swarm

import (
	"swarm/internal/remoter"

	"github.com/anthdm/hollywood/actor"
)

func Run(config SwarmConfig, opts ...Option) error {
	for _, opt := range opts {
		config = opt(config)
	}

	config.adapter = remoter.NewLocalAdapter()

	engine, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(remoter.NewLocalRemoter(config.adapter, "swarm")))
	if err != nil {
		return err
	}

	swarmPID := engine.Spawn(newSwarm(config), "swarm")

	done := make(chan struct{})
	engine.SpawnFunc(func(act *actor.Context) {
		switch act.Message().(type) {
		case actor.Started:
			act.Send(swarmPID, registerListenerEvent{})
		case swarmDoneEvent:
			done <- struct{}{}
		}
	}, "swarm-listener")

	<-done
	return nil
}
