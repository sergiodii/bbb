package channel

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeChannel(t *testing.T) {
	t.Run("Test SafeChannel with int type", func(t *testing.T) {

		// Arrange
		ch := NewSafeChannel[int](1)

		wg := sync.WaitGroup{}
		wg.Add(2)

		// Act & Assert
		go func() {
			defer wg.Done()
			for v := range ch.Receive() {
				assert.Equal(t, 1, v)
			}
		}()

		go func() {
			defer wg.Done()
			ch.Send(1)
			ch.Close() // fecha o canal após enviar
			assert.Equal(t, true, ch.IsClosed())

			err := ch.Send(2) // tenta enviar após o fechamento
			assert.Equal(t, ErrCloseChannel, err)
		}()

		wg.Wait()
	})

	t.Run("Test SafeChannel Close Idempotency", func(t *testing.T) {

		// Arrange
		ch := NewSafeChannel[int](1)

		wg := sync.WaitGroup{}
		wg.Add(2)

		// Act & Assert
		go func() {
			defer wg.Done()
			for v := range ch.Receive() {
				t.Logf("Received: %d", v)
			}
		}()

		go func() {
			defer wg.Done()
			ch.Send(1)
			assert.Equal(t, false, ch.IsClosed())
			ch.Close()
			ch.Close() // tenta fechar novamente
			assert.Equal(t, true, ch.IsClosed())
		}()

		wg.Wait()
	})

	t.Run("Test SafeChannel Receive After Close", func(t *testing.T) {

		// Arrange
		ch := NewSafeChannel[int](1)

		wg := sync.WaitGroup{}
		wg.Add(2)

		// Act
		ch.Send(1)
		ch.Close()
		err := ch.Send(2) // tenta enviar após o fechamento

		// Assert
		assert.Equal(t, ErrCloseChannel, err)
		assert.Equal(t, true, ch.IsClosed())
	})
}
