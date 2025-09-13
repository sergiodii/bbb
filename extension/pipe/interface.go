package pipe

import "context"

type Pipe[T any] interface {
	Enqueue(...func(context.Context, T) (T, error))
	Execute(context.Context, T) (T, error)
}
