package pipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/sergiodii/bbb/extension/slice"

	"golang.org/x/sync/errgroup"
)

var ONF error = errors.New("ObjectNotFound")

type executor[T any] struct {
	tasks []func(context.Context, T) (T, error)
}

func (q *executor[T]) Enqueue(task ...func(context.Context, T) (T, error)) {
	q.tasks = append(q.tasks, task...)
}

func (q *executor[T]) Execute(ctx context.Context, executionType string, target T) (T, error) {
	switch ParseExecutionType(executionType) {
	case SEQUENTIAL:
		return q.executeSequentially(ctx, target)
	case CONCURRENT:
		return q.executeConcurrently(ctx, target)
	case SEQUENTIAL_WITH_FIRST_RESULT:
		return q.executeSequentiallyWithFirstResult(ctx, target)
	case SEQUENTIAL_BLOCKING_ONLY_FIRST:
		return q.executeBlockingFirst(ctx, target)
	default:
		return target, errors.New("unknown execution type: " + executionType)
	}
}

func (q *executor[T]) executeBlockingFirst(ctx context.Context, target T) (T, error) {

	for ind := 0; ind < len(q.tasks); ind++ {
		if ind != 0 {
			select {
			case <-ctx.Done():
				return target, ctx.Err()
			default:
				go func(i int) {
					_, err := q.tasks[i](context.Background(), target)
					if err != nil {
						fmt.Println("Error in non-blocking task:", err)
					}
				}(ind)
			}
		}
	}

	if len(q.tasks) > 0 {
		return q.tasks[0](ctx, target)
	}

	return target, nil
}

func (q *executor[T]) executeConcurrently(ctx context.Context, target T) (T, error) {
	var eg errgroup.Group
	for _, taskList := range slice.TransformSliceToMultipleSlices(q.tasks, 10) {
		// this loop is to limit the number of goroutines running concurrently
		// to avoid overwhelming the system
		// you can adjust the number 10 to a value that makes sense for your use case
		for _, task := range taskList {
			eg.Go(func() error {
				_, err := task(ctx, target)
				if errors.Is(err, ONF) {
					// if the error is ObjectNotFound, continue to the next task
					return nil
				}
				return err
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return target, err
	}

	return target, nil
}

func (q *executor[T]) executeSequentiallyWithFirstResult(ctx context.Context, target T) (T, error) {
	var err error
	for _, task := range q.tasks {
		result, e := task(ctx, target)
		if e == nil {
			return result, nil
		}

		if errors.Is(e, ONF) {

			// if the error is ObjectNotFound, continue to the next task
			continue
		}

		// for other errors, log them
		// the error is stored but the execution continues
		// this error should be logged in a real-world scenario
		errors.Join(err, e)
	}
	return target, err

}

func (q *executor[T]) executeSequentially(ctx context.Context, target T) (T, error) {
	for _, task := range q.tasks {
		t, err := task(ctx, target)
		if err != nil {
			// if the error is ObjectNotFound, continue to the next task
			if errors.Is(err, ONF) {
				continue
			}
			return target, err
		}
		target = t
	}

	return target, nil
}

func NewPipe[T any]() Pipe[T] {
	return &executor[T]{}
}
