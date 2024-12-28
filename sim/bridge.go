package sim

import "github.com/anthdm/hollywood/actor"

type bridgeActor struct {
	simulatorEngine *actor.Engine
	simulatorPID    *actor.PID
}

func newBridgeActor(simulatorEngine *actor.Engine, simulatorPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &bridgeActor{
			simulatorEngine: simulatorEngine,
			simulatorPID:    simulatorPID,
		}
	}
}

func (bridge *bridgeActor) Receive(act *actor.Context) {
	switch msg := act.Message().(type) {
	case actor.Initialized:
	case actor.Started:
		act.Engine().Subscribe(act.PID())
	case actor.Stopped:
		act.Engine().Unsubscribe(act.PID())
	default:
		bridge.simulatorEngine.Send(bridge.simulatorPID, simulationEvent{
			event: msg,
		})
	}
}
