package pipe

import (
	"context"
	"fmt"
)

type sequentiallyBlockingFirstResultPipe[T any] struct {
	tasks []func(context.Context, T) (T, error)
}

// NewSequentiallyBlockingFirstResultPipe creates a new Pipe that executes tasks sequentially
// but only the first task is blocking, the rest are non-blocking
// This means that the first task will be executed and the output of this task will be the input of the next tasks
// The next tasks will be executed in separate goroutines and their results will be ignored
func NewSequentiallyBlockingFirstResultPipe[T any](tasks ...func(context.Context, T) (T, error)) Pipe[T] {
	pipe := &sequentiallyBlockingFirstResultPipe[T]{}

	if len(tasks) == 0 {
		return pipe
	}

	pipe.Enqueue(tasks...)
	return pipe
}
func (p *sequentiallyBlockingFirstResultPipe[T]) Enqueue(tasks ...func(context.Context, T) (T, error)) {
	p.tasks = append(p.tasks, tasks...)
}

func (p *sequentiallyBlockingFirstResultPipe[T]) Execute(ctx context.Context, target T) (T, error) {
	for ind := 0; ind < len(p.tasks); ind++ {
		if ind != 0 {
			select {
			case <-ctx.Done():
				return target, ctx.Err()
			default:
				go func(i int) {
					_, err := p.tasks[i](context.Background(), target)
					if err != nil {
						// log the error, but do not return it
						// in a real-world scenario, you might want to use a logging library
						// here we just print it to the console
						// but in production, i consider creating a channel to collect these errors
						// and return them at the end of the execution
						// or use a logging library to log them
						fmt.Println("Error in non-blocking task:", err)
					}
				}(ind)
			}
		}
	}

	if len(p.tasks) > 0 {
		return p.tasks[0](ctx, target)
	}

	return target, nil
}
