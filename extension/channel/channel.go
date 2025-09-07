package channel

import (
	"errors"
	"sync"
)

var ErrCloseChannel = errors.New("CLOSECHANNEL")

type safeChannel[T any] struct {
	c        chan T
	o        sync.Once
	isClosed bool
	m        sync.RWMutex
}

func NewSafeChannel[T any](bufferSize int) SafeChannel[T] {
	return &safeChannel[T]{
		c:        make(chan T, bufferSize),
		o:        sync.Once{},
		isClosed: false,
	}
}

func (s *safeChannel[T]) Send(value T) error {
	// lock to read the isClosed flag
	s.m.RLock()
	defer s.m.RUnlock()
	if s.isClosed {
		return ErrCloseChannel
	}
	s.c <- value

	return nil
}

func (s *safeChannel[T]) Close() {

	// ensure the channel is closed only once
	s.o.Do(func() {
		close(s.c)

		// lock to set the isClosed flag
		s.m.Lock()
		defer s.m.Unlock()

		s.isClosed = true
	})
}

func (s *safeChannel[T]) Receive() <-chan T {
	return s.c
}

func (s *safeChannel[T]) IsClosed() bool {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.isClosed
}
