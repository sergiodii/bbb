package pipe

import (
	"context"
	"errors"
	"log"
)

// QNF (QueryNotFound) is a sentinel error used to indicate that a requested object was not found.
var ONF error = errors.New("ObjectNotFound")

// NewPipe creates a new Pipe based on the execution type
// It uses the Strategy Pattern to select the appropriate Pipe implementation
// based on the execution type
// The available execution types are:
// - SEQUENTIAL: executes tasks sequentially, where the output of one task is the input of the next task
// - CONCURRENT: executes tasks concurrently, where all tasks receive the same input and their outputs are ignored
// - SEQUENTIAL_WITH_FIRST_RESULT: executes tasks sequentially, but returns the result of the first successful task
// - SEQUENTIAL_BLOCKING_ONLY_FIRST: executes tasks sequentially, but only the first task can modify the input for the next tasks
//
// If an invalid execution type is provided, the function will log a fatal error
// and terminate the program
func NewPipe[T any](executionType ExecutionType, tasks ...func(context.Context, T) (T, error)) Pipe[T] {
	switch executionType {
	case SEQUENTIAL:
		return NewSequentiallyPipe(tasks...)
	case CONCURRENT:
		return NewConcurrentlyPipe(tasks...)
	case SEQUENTIAL_WITH_FIRST_RESULT:
		return NewSequentiallyWithFirstResultPipe(tasks...)
	case SEQUENTIAL_BLOCKING_ONLY_FIRST:
		return NewSequentiallyBlockingFirstResultPipe(tasks...)
	default:
		log.Fatalln("Invalid execution type")
		return nil
	}
}
