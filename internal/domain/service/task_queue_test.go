package service_test

import (
	"testing"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/service"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

type (
	// delayedFactory creates a task that simulates a delay before completion.
	delayedFactory struct{}

	// delayedTask sleeps for a short duration to mimic work.
	delayedTask struct{}
)

// New returns a task that sleeps for 200ms.
func (*delayedFactory) New(_ *model.Task) task.ExecutableTask {
	return &delayedTask{}
}

// Run simulates task execution by sleeping for 200ms.
func (*delayedTask) Run() error {
	time.Sleep(200 * time.Millisecond)
	return nil
}

// waitUntilDone polls the task until it completes or times out.
func waitUntilDone(t *testing.T, manager *service.TaskManager, id string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)

	for {
		tsk, ok := manager.GetTask(id)

		if !ok {
			t.Fatal("task not found")
		}
		if tsk.Status == model.TaskStatusDone || tsk.Status == model.TaskStatusFailed {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("task %s did not complete in time", id)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

// TestSequentialExecution_PerTaskType verifies that tasks of the same type
// are executed sequentially - not in parallel.
func TestSequentialExecution_PerTaskType(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("delayed", &delayedFactory{})

	start := time.Now()

	t1 := manager.CreateTask("delayed")
	t2 := manager.CreateTask("delayed")

	waitUntilDone(t, manager, t1.ID)
	waitUntilDone(t, manager, t2.ID)

	if elapsed := time.Since(start); elapsed < 400*time.Millisecond {
		t.Errorf("expected sequential execution, but took only %v", elapsed)
	}
}
