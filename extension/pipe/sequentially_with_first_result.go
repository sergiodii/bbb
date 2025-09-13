package pipe

import (
	"context"
	"errors"
)

type sequentiallyWithFirstResultPipe[T any] struct {
	tasks []func(context.Context, T) (T, error)
}

// NewSequentiallyWithFirstResultPipe creates a new Pipe that executes tasks sequentiall
// but returns the result of the first successful task
// this means that the first task that succeeds will return its result immediately
// and the rest of the tasks will not be executed
// If a task returns an ONF (ObjectNotFound) error, the execution continues to the next task
// If all tasks return an error, the errors is joined and returned
func NewSequentiallyWithFirstResultPipe[T any](tasks ...func(context.Context, T) (T, error)) Pipe[T] {
	pipe := &sequentiallyWithFirstResultPipe[T]{}

	if len(tasks) == 0 {
		return pipe
	}

	pipe.Enqueue(tasks...)
	return pipe
}
func (p *sequentiallyWithFirstResultPipe[T]) Enqueue(tasks ...func(context.Context, T) (T, error)) {
	p.tasks = append(p.tasks, tasks...)
}

func (p *sequentiallyWithFirstResultPipe[T]) Execute(ctx context.Context, target T) (T, error) {
	var err error
	for _, task := range p.tasks {
		result, e := task(ctx, target)
		if e == nil {
			return result, nil
		}

		// if the error is ObjectNotFound, continue to the next task
		if errors.Is(e, ONF) {
			continue
		}

		// for other errors, log them
		// the error is stored but the execution continues
		// this error should be logged in a real-world scenario
		errors.Join(err, e)
	}
	return target, err
}
