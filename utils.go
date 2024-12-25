package swarm

import (
	"time"

	"github.com/anthdm/hollywood/actor"
)

func sendWithDelay(act *actor.Context, pid *actor.PID, msg any, delay time.Duration) {
	sender := act.PID()
	cancel := make(chan struct{})
	act.SpawnChildFunc(func(act *actor.Context) {
		switch act.Message().(type) {
		case actor.Started:
			go func() {
				timer := time.NewTimer(delay)
				select {
				case <-timer.C:
					act.Engine().SendWithSender(pid, msg, sender)
				case <-cancel:
					timer.Stop()
				}
			}()
		case actor.Stopped:
			close(cancel)
		}
	}, "delayed-sender")
}
