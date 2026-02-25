package main

type Task func() error

type TaskRunner struct {
	tasks []Task
}

func (tr *TaskRunner) AddTask(t Task) {
	tr.tasks = append(tr.tasks, t)
}

func (tr *TaskRunner) RunAll() []error {
	errorList := []error{}
	for _, task := range tr.tasks {
		if err := task(); err != nil {
			errorList = append(errorList, err)
		}

	}
	return errorList
}
