package vote

import "context"

type ExecutionType interface {
	String() string
}

type Pipe[T any] interface {
	Enqueue(...func(context.Context, T) (T, error))
	Execute(context.Context, string, T) (T, error)
}
