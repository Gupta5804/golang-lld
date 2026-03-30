package engine

import (
	"sync"

	"github.com/Gupta5804/golang-lld/systems/phase1/05_admin_task_engine/task"
)

type ExecutionEngine struct {
	history []task.Task
	mu      sync.RWMutex
}

func (e *ExecutionEngine) Execute(t task.Task) error {
	if err := t.Execute(); err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.history = append(e.history,t)
	return nil
}

func (e *ExecutionEngine) UndoLast() error {
	e.mu.Lock()
	
	if len(e.history) == 0 {
		return nil
	}
	
	lastIdx := len(e.history) - 1
	lastTask := e.history[lastIdx]
	e.history = e.history[:lastIdx]
	
	e.mu.Unlock()
	return lastTask.Undo()
}
