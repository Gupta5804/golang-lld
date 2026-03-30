package task_test

import (
	"errors"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/05_admin_task_engine/task"
)

var ErrMockFailure = errors.New("mock execution failed")

type MockTask struct {
	ExecuteCalled bool
	UndoCalled    bool
	ShouldFail    bool
}

func (m *MockTask) Execute() error {
	m.ExecuteCalled = true
	if m.ShouldFail {
		return ErrMockFailure
	}
	return nil
}
func (m *MockTask) Undo() error {
	m.UndoCalled = true
	return nil
}

func TestBundleTask_ExecutionAndRollback(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() (*task.BundleTask, []*MockTask)
		expectedError error
		assertMocks   func(t *testing.T, mocks []*MockTask)
		triggerUndo   bool
	}{
		{
			name: "Happy Path (All tasks succeed)",
			setup: func() (*task.BundleTask, []*MockTask) {
				task1 := &MockTask{}
				task2 := &MockTask{}
				bundle := task.NewBundleTask(task1, task2)
				return bundle, []*MockTask{task1, task2}
			},
			expectedError: nil,
			assertMocks: func(t *testing.T, mocks []*MockTask) {
				for i, mock := range mocks {
					if !mock.ExecuteCalled {
						t.Errorf("Mock %d:Expected Execute to be called but was not called", i)
					}
					if mock.UndoCalled {
						t.Errorf("Mock %d:Expected Undo to be not called, but was called", i)
					}
				}
			},
			triggerUndo: false,
		},
		{
			name: "Partial Failure(middle task fails, previous tasks must undo)",
			setup: func() (*task.BundleTask, []*MockTask) {
				task1 := &MockTask{}
				task2 := &MockTask{ShouldFail: true}
				task3 := &MockTask{}

				bundle := task.NewBundleTask(task1, task2, task3)
				return bundle, []*MockTask{task1, task2, task3}
			},
			expectedError: ErrMockFailure,
			assertMocks: func(t *testing.T, mocks []*MockTask) {
				mock1, mock2, mock3 := mocks[0], mocks[1], mocks[2]

				if !mock1.ExecuteCalled {
					t.Errorf("Task1: Expected mock to call execute but didnt")
				}
				if !mock2.ExecuteCalled {
					t.Errorf("Task2: Expected mock to call Execute but didnt")
				}
				if mock3.ExecuteCalled {
					t.Errorf("Task3: Expected mock to not call Execute, but it did")
				}
				if !mock1.UndoCalled {
					t.Errorf("Task1: Expected mock to call Undo, but did not")
				}
				if mock2.UndoCalled {
					t.Errorf("Task2: Expected mock to not call Undo, but it did")
				}
				if mock3.UndoCalled {
					t.Errorf("Task3: Expected mock to not call Undo, but it did")
				}

			},
			triggerUndo: false,
		},
		{
			name: "Explicit Undo(All Tasks successfully undo)",
			setup: func() (*task.BundleTask, []*MockTask) {
				task1 := &MockTask{}
				task2 := &MockTask{}
				bundle := task.NewBundleTask(task1, task2)
				return bundle, []*MockTask{task1, task2}
			},
			expectedError: nil,
			assertMocks: func(t *testing.T, mocks []*MockTask) {
				for i, mock := range mocks {
					if !mock.ExecuteCalled {
						t.Errorf("Mock %d: Expected Mock to call execute but didnt", i)
					}
					if !mock.UndoCalled {
						t.Errorf("Mock %d: Expected mock to call Undo but didnt", i)
					}
				}
			},
			triggerUndo: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bundle, mocks := tc.setup()
			err := bundle.Execute()
			if tc.triggerUndo {
				err = bundle.Undo()
			}
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("Expected error %v, got %v", tc.expectedError, err)
			}
			tc.assertMocks(t, mocks)
		})
	}
}
