package server

import (
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/swarm/sim"
)

func NewSwarmHandler(engine *actor.Engine, pid *actor.PID) http.Handler {
	server := swarmServer{
		engine: engine,
		pid:    pid,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.getRootHandler)
	mux.HandleFunc("POST /start", server.postStartHandler)
	mux.HandleFunc("POST /stop", server.postStopHandler)

	return mux
}

type swarmServer struct {
	engine *actor.Engine
	pid    *actor.PID
}

func (server swarmServer) getRootHandler(w http.ResponseWriter, r *http.Request) {
}

func (server swarmServer) postStartHandler(w http.ResponseWriter, r *http.Request) {
	server.engine.Send(server.pid, sim.Start{})
	w.WriteHeader(http.StatusOK)
}

func (server swarmServer) postStopHandler(w http.ResponseWriter, r *http.Request) {
	server.engine.Send(server.pid, sim.Stop{})
	w.WriteHeader(http.StatusOK)
}
