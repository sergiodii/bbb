package vote

import "context"

// Outbound Ports
type Pipe[T any] interface {
	Enqueue(...func(context.Context, T) (T, error))
	Execute(context.Context, T) (T, error)
}
