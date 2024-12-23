package swarm

type startSimulationEvent struct{}

type simulationDoneEvent struct {
	seed int64
}

type simulationErrorEvent struct {
	seed int64
	err  error
}

type registerListenerEvent struct{}

type swarmDoneEvent struct{}
