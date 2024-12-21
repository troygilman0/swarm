package remoter

import (
	"fmt"
	"sync"

	"github.com/anthdm/hollywood/actor"
)

type Adapter interface {
	start(address string, engine *actor.Engine) error
	stop(address string) *sync.WaitGroup
	send(pid *actor.PID, msg any, sender *actor.PID)
}

type localAdapter struct {
	engines map[string]*actor.Engine
	lock    sync.RWMutex
}

func NewLocalAdapter() Adapter {
	return &localAdapter{
		engines: make(map[string]*actor.Engine),
	}
}

func (adapter *localAdapter) start(address string, engine *actor.Engine) error {
	adapter.lock.Lock()
	defer adapter.lock.Unlock()
	if _, ok := adapter.engines[address]; ok {
		return fmt.Errorf("address already taken: %s", address)
	}
	adapter.engines[address] = engine
	return nil
}

func (adapter *localAdapter) stop(address string) *sync.WaitGroup {
	adapter.lock.Lock()
	defer adapter.lock.Unlock()
	if _, ok := adapter.engines[address]; !ok {
		return &sync.WaitGroup{}
	}
	delete(adapter.engines, address)
	return &sync.WaitGroup{}
}

func (adapter *localAdapter) send(pid *actor.PID, msg any, sender *actor.PID) {
	adapter.lock.RLock()
	defer adapter.lock.RUnlock()
	address := pid.GetAddress()
	engine, ok := adapter.engines[address]
	if !ok {
		return
	}
	engine.SendWithSender(pid, msg, sender)
}
