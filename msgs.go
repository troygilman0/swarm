package swarm

type StartSimulation struct{}

type SimulationDoneEvent struct {
	Seed int64
}

type SimulationErrorEvent struct {
	Seed  int64
	Error error
}

type RegisterListener struct{}

type SwarmDoneEvent struct{}
