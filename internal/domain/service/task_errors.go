package service

import "errors"

// Predefined errors returned by the TaskManager methods.
var (
	ErrTaskNotFound          = errors.New("task not found")
	ErrTaskInProgress        = errors.New("task in progress")
	ErrTaskAlreadyExists     = errors.New("task already exists")
	ErrTaskQueueLimitReached = errors.New("task queue limit reached")
	ErrTaskUnknownType       = errors.New("task unknown type")
)
