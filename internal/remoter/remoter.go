package remoter

import (
	"sync"

	"github.com/anthdm/hollywood/actor"
)

type localRemoter struct {
	address string
	adapter Adapter
}

func NewLocalRemoter(adapter Adapter, address string) actor.Remoter {
	return &localRemoter{
		address: address,
		adapter: adapter,
	}
}

func (remoter *localRemoter) Address() string {
	return remoter.address
}

func (remoter *localRemoter) Send(pid *actor.PID, msg any, sender *actor.PID) {
	remoter.adapter.send(pid, msg, sender)
}

func (remoter *localRemoter) Start(engine *actor.Engine) error {
	return remoter.adapter.start(remoter.address, engine)
}

func (remoter *localRemoter) Stop() *sync.WaitGroup {
	return remoter.adapter.stop(remoter.address)
}
