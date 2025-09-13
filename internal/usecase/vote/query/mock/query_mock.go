package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type PipeMock[T any] struct {
	mock.Mock
}

func (p *PipeMock[T]) Enqueue(fns ...func(context.Context, T) (T, error)) {
	p.Called(fns)
}
func (p *PipeMock[T]) Execute(ctx context.Context, dto T) (T, error) {
	args := p.Called(ctx, dto)
	return args.Get(0).(T), args.Error(1)
}

func NewPipeMock[T any]() *PipeMock[T] {
	return &PipeMock[T]{}
}
