package notifier_test

import (
	"sync/atomic"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/notifier"
)

type MockCommand struct {
	ExecuteFunc func() error
}

func (m *MockCommand) Execute() error {
	return m.ExecuteFunc()
}

func TestDispatcherHub_WorkerPool(t *testing.T) {
	commandChan := make(chan domain.Command, 10) // 10 buffer
	hub := notifier.NewDispatcherHub(commandChan, 3) // 3 workers
	var counter atomic.Int32
	for i := 0; i < 5; i++ {
		cmd := &MockCommand{
			ExecuteFunc: func()error{
				counter.Add(1)
				return nil
			},
		}
		commandChan <- cmd
	}
	hub.Start()
	close(commandChan)
	hub.Stop()
	if actual := counter.Load();actual != 5{
		t.Errorf("Expected 5 commands got %d",actual)
	}
}
