package channel

type SafeChannel[T any] interface {
	Send(value T) error
	Close()
	Receive() <-chan T
	IsClosed() bool
}
