package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman0/swarm/sim"
)

func NewSwarmHandler(initializer actor.Producer, opts ...sim.Option) http.Handler {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	pid := engine.Spawn(newControllerActor(initializer, opts...), "controller")

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
	buildDir := "server/ui/build/client" // Path to Remix build output
	remixPath := filepath.Join(buildDir, r.URL.Path)

	if _, err := os.Stat(remixPath); os.IsNotExist(err) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, remixPath)
}

func (server swarmServer) postStartHandler(w http.ResponseWriter, r *http.Request) {
	server.engine.Send(server.pid, start{})
	w.WriteHeader(http.StatusOK)
}

func (server swarmServer) postStopHandler(w http.ResponseWriter, r *http.Request) {
	server.engine.Send(server.pid, stop{})
	w.WriteHeader(http.StatusOK)
}
