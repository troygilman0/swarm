package remoter

import (
	"sync"

	"github.com/anthdm/hollywood/actor"
)

type remoter struct {
	address string
	adapter Adapter
}

func NewRemoter(adapter Adapter, address string) actor.Remoter {
	return &remoter{
		address: address,
		adapter: adapter,
	}
}

func (remoter *remoter) Address() string {
	return remoter.address
}

func (remoter *remoter) Send(pid *actor.PID, msg any, sender *actor.PID) {
	remoter.adapter.Send(pid, msg, sender)
}

func (remoter *remoter) Start(engine *actor.Engine) error {
	return remoter.adapter.Start(remoter.address, engine)
}

func (remoter *remoter) Stop() *sync.WaitGroup {
	return remoter.adapter.Stop(remoter.address)
}
