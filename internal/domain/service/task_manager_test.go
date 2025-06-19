package service_test

import (
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

	tsk := manager.CreateTask("mock")

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
// with a 'failed' status if the task type is not registered.
func TestCreateTask_UnknownType(t *testing.T) {
	manager := service.NewTaskManager()
	tsk := manager.CreateTask("unknown")

	if tsk.Status != model.TaskStatusFailed {
		t.Errorf("expected status 'failed', got %q", tsk.Status)
	}
	if tsk.Result == "" {
		t.Error("expected error message in task result")
	}
}

// TestGetTask_Found checks that a created task can be retrieved.
func TestGetTask_Found(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})

	tsk := manager.CreateTask("mock")
	found, ok := manager.GetTask(tsk.ID)

	if !ok {
		t.Fatal("expected task to be found")
	}
	if found.ID != tsk.ID {
		t.Error("returned task ID mismatch")
	}
}

// TestGetTask_NotFound checks that retrieving a non-existent task fails.
func TestGetTask_NotFound(t *testing.T) {
	manager := service.NewTaskManager()
	_, ok := manager.GetTask("non-existent")

	if ok {
		t.Error("expected not to find nonexistent task")
	}
}

// TestDeleteTask verifies that a finished task can be deleted successfully.
func TestDeleteTask(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})

	tsk := manager.CreateTask("mock")
	tsk.Status = model.TaskStatusDone
	deleted, locked := manager.DeleteTask(tsk.ID)

	if !deleted {
		t.Error("expected task to be deleted")
	}
	if locked {
		t.Error("did not expect task to be locked")
	}
	if _, ok := manager.GetTask(tsk.ID); ok {
		t.Error("expected task to be gone")
	}
}

// TestDeleteTask_Running checks that running tasks cannot be deleted.
func TestDeleteTask_Running(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("mock", &mockFactory{})

	tsk := manager.CreateTask("mock")
	tsk.Status = model.TaskStatusRunning
	deleted, locked := manager.DeleteTask(tsk.ID)

	if deleted {
		t.Error("should not delete a running task")
	}
	if !locked {
		t.Error("expected task to be locked")
	}
}
