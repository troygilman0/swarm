package server

import (
	"net/http"

	"github.com/anthdm/hollywood/actor"
)

func sendMessageHandler(engine *actor.Engine, pid *actor.PID, msg any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
