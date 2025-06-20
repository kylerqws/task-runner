package service_test

import (
	"errors"
	"testing"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/service"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

type (
	mockFactory struct{} // mockFactory returns a mock task that does nothing but succeeds.
	mockTask    struct{} // mockTask is a dummy task that always succeeds.
)

// New returns a mock task.
func (*mockFactory) New(_ *model.Task) task.ExecutableTask {
	return &mockTask{}
}

// Run simulates success.
func (*mockTask) Run() error {
	return nil
}

// TestCreateTask_Success checks that a task is created properly with a known type.
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
	if tsk.Type != "mock" {
		t.Errorf("expected type 'mock', got %q", tsk.Type)
	}
	if tsk.Status != model.TaskStatusPending {
		t.Errorf("expected status 'pending', got %q", tsk.Status)
	}
	if tsk.ID == "" {
		t.Error("expected non-empty task ID")
	}
}

// TestCreateTask_UnknownType ensures an error is returned for an unregistered type.
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

// TestCreateTask_EmptyType verifies that creating a task without type fails.
func TestCreateTask_EmptyType(t *testing.T) {
	manager := service.NewTaskManager()
	_, err := manager.CreateTask("")
	if err == nil {
		t.Fatal("expected error for empty task type")
	}
	if !errors.Is(err, service.ErrTaskUnknownType) {
		t.Errorf("expected ErrTaskUnknownType, got %v", err)
	}
}

// TestGetTask_Found ensures that a task can be retrieved by ID.
func TestGetTask_Found(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})
	tsk, _ := manager.CreateTask("mock")
	found, err := manager.GetTask(tsk.ID)
	if err != nil {
		t.Fatalf("expected to find task, got error: %v", err)
	}
	if found.ID != tsk.ID {
		t.Error("returned task ID mismatch")
	}
}

// TestGetTask_NotFound checks that retrieving an unknown task returns an error.
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

// TestDeleteTask ensures that a completed task can be deleted.
func TestDeleteTask(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})
	tsk, _ := manager.CreateTask("mock")
	tsk.Status = model.TaskStatusDone

	err := manager.DeleteTask(tsk.ID)
	if err != nil {
		t.Errorf("expected task to be deleted, got error: %v", err)
	}
	_, err = manager.GetTask(tsk.ID)
	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Error("expected task to be gone")
	}
}

// TestDeleteTask_Running verifies that running tasks cannot be deleted.
func TestDeleteTask_Running(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})
	tsk, _ := manager.CreateTask("mock")
	tsk.Status = model.TaskStatusRunning

	err := manager.DeleteTask(tsk.ID)
	if err == nil {
		t.Fatal("expected error for running task")
	}
	if !errors.Is(err, service.ErrTaskInProgress) {
		t.Errorf("expected ErrTaskInProgress, got %v", err)
	}
}
