package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/service"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

type (
	delayedFactory struct{} // delayedFactory creates a task that completes after a short delay.
	delayedTask    struct{} // delayedTask sleeps for a brief period to simulate work.
)

type (
	blockingFactory struct{}                     // blockingFactory creates a task that never completes.
	blockingTask    struct{ hold chan struct{} } // blockingTask is a task that blocks indefinitely.
)

// New returns a delayed task that sleeps for a short duration.
func (*delayedFactory) New(_ *model.Task) task.ExecutableTask {
	return &delayedTask{}
}

// Run sleeps for 200ms to simulate work.
func (*delayedTask) Run() error {
	time.Sleep(200 * time.Millisecond)
	return nil
}

// New returns a blocking task that never completes by default.
func (*blockingFactory) New(_ *model.Task) task.ExecutableTask {
	return &blockingTask{hold: make(chan struct{})}
}

// Run blocks until the internal channel is closed.
func (b *blockingTask) Run() error {
	<-b.hold
	return nil
}

// waitUntilDone waits for a task to complete or fails on timeout.
func waitUntilDone(t *testing.T, manager *service.TaskManager, id string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		tsk, err := manager.GetTask(id)
		if err != nil {
			t.Fatalf("task not found: %v", err)
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

// TestSequentialExecution_PerTaskType ensures only one task runs at a time per type.
func TestSequentialExecution_PerTaskType(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("delayed", &delayedFactory{})

	start := time.Now()

	t1, err := manager.CreateTask("delayed")
	if err != nil {
		t.Fatalf("unexpected error creating first task: %v", err)
	}
	t2, err := manager.CreateTask("delayed")
	if err != nil {
		t.Fatalf("unexpected error creating second task: %v", err)
	}

	waitUntilDone(t, manager, t1.ID)
	waitUntilDone(t, manager, t2.ID)

	if elapsed := time.Since(start); elapsed < 400*time.Millisecond {
		t.Errorf("expected sequential execution, but took only %v", elapsed)
	}
}

// TestCreateTask_QueueOverflow ensures queue size is enforced.
func TestCreateTask_QueueOverflow(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("blocked", &blockingFactory{})

	for i := 0; i < 100; i++ {
		if _, err := manager.CreateTask("blocked"); err != nil {
			t.Fatalf("unexpected error while filling queue: %v", err)
		}
	}

	tk, err := manager.CreateTask("blocked")
	if err == nil {
		t.Fatal("expected error for queue overflow, got nil")
	}
	if !errors.Is(err, service.ErrTaskQueueLimitReached) {
		t.Errorf("expected ErrTaskQueueLimitReached, got: %v", err)
	}
	if tk != nil {
		t.Error("expected nil task on overflow")
	}
}
