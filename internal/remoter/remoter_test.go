package remoter

import (
	"testing"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func TestSendMessage(t *testing.T) {
	adapter := NewLocalAdapter()

	engine1, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(NewRemoter(adapter, "engine-1")))
	if err != nil {
		t.Error(err)
	}

	engine2, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(NewRemoter(adapter, "engine-2")))
	if err != nil {
		t.Error(err)
	}

	result := make(chan string)
	pid := engine2.SpawnFunc(func(act *actor.Context) {
		switch msg := act.Message().(type) {
		case string:
			result <- msg
		}
	}, "tester")

	engine1.Send(pid, "hello world")

	select {
	case msg := <-result:
		if msg != "hello world" {
			t.Errorf("result msg is incorrect")
		}
	case <-time.After(time.Second):
		t.Errorf("timed out")
	}
}
