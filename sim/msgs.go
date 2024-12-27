package sim

type Start struct{}

type Stop struct{}

type startSimulation struct{}

type stopSimulation struct{}

type SimulationDoneEvent struct {
	Seed int64
}

type SimulationErrorEvent struct {
	Seed  int64
	Error error
}

type RegisterListener struct{}

type ManagerDoneEvent struct{}
