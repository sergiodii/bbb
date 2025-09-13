package pipe

import (
	"context"
	"errors"
)

// sequentiallyPipe  applies the Strategy Pattern
type sequentiallyPipe[T any] struct {
	tasks []func(context.Context, T) (T, error)
}

// NewSequentiallyPipe creates a new Pipe that executes tasks sequentially
// Sequentially means that each task will be executed one after the other
// and the output of one task will be the input of the next task
func NewSequentiallyPipe[T any](tasks ...func(context.Context, T) (T, error)) Pipe[T] {
	pipe := &sequentiallyPipe[T]{}

	if len(tasks) == 0 {
		return pipe
	}

	pipe.Enqueue(tasks...)
	return pipe
}

func (p *sequentiallyPipe[T]) Enqueue(tasks ...func(context.Context, T) (T, error)) {
	p.tasks = append(p.tasks, tasks...)
}

func (p *sequentiallyPipe[T]) Execute(ctx context.Context, input T) (T, error) {
	for _, task := range p.tasks {
		t, err := task(ctx, input)
		if err != nil {

			// if the error is ObjectNotFound, continue to the next task
			if errors.Is(err, ONF) {
				continue
			}
			return input, err
		}
		input = t
	}

	return input, nil
}
