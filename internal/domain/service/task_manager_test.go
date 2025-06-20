package service_test

import (
	"errors"
	"testing"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/service"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

type (
	// mockFactory returns a mock task that does nothing but succeeds.
	mockFactory struct{}

	// mockTask is a dummy task that always succeeds.
	mockTask struct{}
)

// New returns a new instance of a mock task that always succeeds.
func (*mockFactory) New(_ *model.Task) task.ExecutableTask {
	return &mockTask{}
}

// Run is a no-op implementation that always returns nil.
func (*mockTask) Run() error {
	return nil
}

// TestCreateTask_Success verifies that a task is created successfully
// when a valid task type is registered.
func TestCreateTask_Success(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})

	tsk, err := manager.CreateTask("mock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tsk == nil {
		t.Fatal("expected task to be created")
	}
	if tsk.Status != model.TaskStatusPending {
		t.Errorf("expected status 'pending', got %q", tsk.Status)
	}
	if tsk.ID == "" {
		t.Error("expected non-empty task ID")
	}
}

// TestCreateTask_UnknownType ensures that task creation fails
// if the task type is not registered.
func TestCreateTask_UnknownType(t *testing.T) {
	manager := service.NewTaskManager()
	tsk, err := manager.CreateTask("unknown")

	if err == nil {
		t.Fatal("expected error for unknown task type")
	}
	if !errors.Is(err, service.ErrTaskUnknownType) {
		t.Errorf("expected ErrTaskUnknownType, got %v", err)
	}
	if tsk != nil {
		t.Error("expected returned task to be nil")
	}
}

// TestGetTask_Found checks that a created task can be retrieved.
func TestGetTask_Found(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})

	tsk, err := manager.CreateTask("mock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := manager.GetTask(tsk.ID)
	if err != nil {
		t.Fatalf("expected to find task, got error: %v", err)
	}

	if found.ID != tsk.ID {
		t.Error("returned task ID mismatch")
	}
}

// TestGetTask_NotFound checks that retrieving a non-existent task returns an error.
func TestGetTask_NotFound(t *testing.T) {
	manager := service.NewTaskManager()
	_, err := manager.GetTask("non-existent")

	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

// TestDeleteTask verifies that a finished task can be deleted successfully.
func TestDeleteTask(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})

	tsk, err := manager.CreateTask("mock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tsk.Status = model.TaskStatusDone

	err = manager.DeleteTask(tsk.ID)
	if err != nil {
		t.Errorf("expected task to be deleted, got error: %v", err)
	}

	_, err = manager.GetTask(tsk.ID)
	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Error("expected task to be deleted from manager")
	}
}

// TestDeleteTask_Running checks that running tasks cannot be deleted.
func TestDeleteTask_Running(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})

	tsk, err := manager.CreateTask("mock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tsk.Status = model.TaskStatusRunning

	err = manager.DeleteTask(tsk.ID)
	if err == nil {
		t.Fatal("expected error when deleting running task")
	}
	if !errors.Is(err, service.ErrTaskIsProgress) {
		t.Errorf("expected ErrTaskIsProgress, got %v", err)
	}
}
