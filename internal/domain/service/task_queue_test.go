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
	// delayedFactory creates a task that simulates a delay before completion.
	delayedFactory struct{}

	// delayedTask sleeps for a short duration to mimic work.
	delayedTask struct{}
)

type (
	// blockingFactory creates tasks that block forever (for queue overflow testing).
	blockingFactory struct{}

	// blockingTask never completes until externally signaled.
	blockingTask struct {
		hold chan struct{}
	}
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

// New returns a blocking task with an unclosed hold channel.
func (*blockingFactory) New(_ *model.Task) task.ExecutableTask {
	return &blockingTask{hold: make(chan struct{})}
}

// Run blocks until the hold channel is closed (never, by default).
func (b *blockingTask) Run() error {
	<-b.hold
	return nil
}

// waitUntilDone polls the task until it completes or times out.
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

// TestSequentialExecution_PerTaskType verifies that tasks of the same type
// are executed sequentially (not in parallel).
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

// TestCreateTask_QueueOverflow ensures task creation fails
// if the queue for the given type exceeds its capacity.
func TestCreateTask_QueueOverflow(t *testing.T) {
	manager := service.NewTaskManager()
	manager.RegisterFactory("blocked", &blockingFactory{})

	// Fill the queue to capacity
	for i := 0; i < 100; i++ {
		_, err := manager.CreateTask("blocked")
		if err != nil {
			t.Fatalf("unexpected error while filling queue: %v", err)
		}
	}

	// Next task should exceed limit
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
