package main

import (
	"fmt"

	"github.com/Gupta5804/golang-lld/systems/phase1/05_admin_task_engine/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/05_admin_task_engine/engine"
	"github.com/Gupta5804/golang-lld/systems/phase1/05_admin_task_engine/task"
)

type DummyRepo struct{}

func (d *DummyRepo) GetByID(id string) (*domain.User, error) {
	return &domain.User{ID: id}, nil
}
func (d *DummyRepo) Save(u *domain.User) error {
	return nil
}

func main() {
	fmt.Println("===== Starting Admin Task Engine =====")

	repo := &DummyRepo{}
	execEngine := &engine.ExecutionEngine{}

	// scenario 1 : a single atomic task
	fmt.Println("\n--- Scenario 1: Single Task ---")
	suspendTask := task.NewSuspendUserTask("user_123", repo)
	execEngine.Execute(suspendTask)
	execEngine.UndoLast()

	// scenario 2 : A bundle task
	fmt.Println("\n--- Scenario 2: Bundle Task ---")
	resetTask := task.NewResetPasswordTask("user_456",repo, "new_secure_password")
	notifyTask := task.NewNotifyTask("user_456","You password is reset",repo)
	bundle := task.NewBundleTask(resetTask,notifyTask)

	execEngine.Execute(bundle)

	fmt.Println("\n --- Admin is rolling back ---")
	execEngine.UndoLast()
}
