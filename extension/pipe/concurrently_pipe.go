package pipe

import (
	"context"
	"errors"

	"github.com/sergiodii/bbb/extension/slice"
	"golang.org/x/sync/errgroup"
)

type concurrentlyPipe[T any] struct {
	tasks []func(context.Context, T) (T, error)
}

// NewConcurrentlyPipe creates a new Pipe that executes tasks concurrently
// Concurrently means that all tasks will be executed at the same time
// and the output of one task will not be the input of the next task
// The input of all tasks will be the same
func NewConcurrentlyPipe[T any](tasks ...func(context.Context, T) (T, error)) Pipe[T] {
	pipe := &concurrentlyPipe[T]{}

	if len(tasks) == 0 {
		return pipe
	}

	pipe.Enqueue(tasks...)
	return pipe
}

func (p *concurrentlyPipe[T]) Enqueue(tasks ...func(context.Context, T) (T, error)) {
	p.tasks = append(p.tasks, tasks...)
}

func (p *concurrentlyPipe[T]) Execute(ctx context.Context, input T) (T, error) {
	var eg errgroup.Group
	for _, taskList := range slice.TransformSliceToMultipleSlices(p.tasks, 10) {
		// this loop is to limit the number of goroutines running concurrently
		// to avoid overwhelming the system
		// you can adjust the number 10 to a value that makes sense for your use case
		for _, task := range taskList {
			eg.Go(func() error {
				_, err := task(ctx, input)

				// if the error is ObjectNotFound, continue to the next task
				if errors.Is(err, ONF) {
					return nil
				}
				return err
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return input, err
	}

	return input, nil
}
