package task

import "github.com/Gupta5804/golang-lld/systems/phase1/05_admin_task_engine/domain"

type Task interface { // command interface
	Execute() error
	Undo() error
}
type UserRepository interface {
	GetByID(id string) (*domain.User, error)
	Save(u *domain.User) error
}

// SuspendUserTask commmand
type SuspendUserTask struct {
	userID        string
	repo          UserRepository
	previousState bool
}

func (t *SuspendUserTask) Execute() error {
	return nil
}
func (t *SuspendUserTask) Undo() error {
	return nil
}
func NewSuspendUserTask(userID string, repo UserRepository) *SuspendUserTask {
	return &SuspendUserTask{
		userID: userID,
		repo:   repo,
	}
}

// ResetPasswordTask Command
type ResetPasswordTask struct {
	userID       string
	repo         UserRepository
	previousPass string
	newPass      string
}

func (t *ResetPasswordTask) Execute() error {
	return nil
}
func (t *ResetPasswordTask) Undo() error {
	return nil
}
func NewResetPasswordTask(userID string, repo UserRepository, newPassword string) *ResetPasswordTask {
	return &ResetPasswordTask{
		userID:  userID,
		repo:    repo,
		newPass: newPassword,
	}
}

// NotifyTask Command
type NotifyTask struct {
	userID  string
	repo    UserRepository
	message string
}

func (t *NotifyTask) Execute() error {
	return nil
}
func (t *NotifyTask) Undo() error {
	return nil
}
func NewNotifyTask(userID string, message string, repo UserRepository) *NotifyTask {
	return &NotifyTask{
		userID:  userID,
		repo:    repo,
		message: message,
	}
}

// BundleTask , composite of commands
type BundleTask struct {
	tasks []Task
}

func (t *BundleTask) Execute() error {
	for i, task := range t.tasks {
		if err := task.Execute(); err != nil {
			for j := i - 1; j >= 0; j-- {
				_ = t.tasks[j].Undo()
			}
			return err
		}
	}
	return nil
}
func (t *BundleTask) Undo() error {
	for i := len(t.tasks) - 1; i >= 0; i-- {
		if err := t.tasks[i].Undo(); err != nil {
			return err
		}
	}
	return nil
}
func NewBundleTask(tasks ...Task) *BundleTask {
	return &BundleTask{
		tasks: tasks,
	}
}
