package server

import (
	"log"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/swarm/sim"
)

type controllerActor struct {
	managerPID  *actor.PID
	initializer actor.Producer
	opts        []sim.Option
}

func newControllerActor(initializer actor.Producer, opts ...sim.Option) actor.Producer {
	return func() actor.Receiver {
		return &controllerActor{
			initializer: initializer,
			opts:        opts,
		}
	}
}

func (controller *controllerActor) Receive(act *actor.Context) {
	log.Printf("%s : %T - %+v\n", act.PID().String(), act.Message(), act.Message())

	switch act.Message().(type) {
	case start:
		controller.managerPID = act.SpawnChild(sim.NewManager(controller.initializer, controller.opts...), "manager")
		act.Send(controller.managerPID, sim.RegisterListener{})

	case stop:
		if controller.managerPID == nil {
			return
		}
		act.Engine().Stop(controller.managerPID).Wait()
		controller.managerPID = nil
	}
}
