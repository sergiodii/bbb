package pipe

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	t.Run("Sequential Execution", func(t *testing.T) {
		// Arrange
		pipe := NewPipe[int]()
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return i + 1, nil
		})
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return i * 2, nil
		})

		// Act
		result, err := pipe.Execute(context.Background(), "SEQUENTIAL", 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Assert

		// Validate that the result is (1 + 1) * 2 = 4
		// The first function adds 1 to the input (1 + 1 = 2)
		// The second function multiplies the result by 2 (2 * 2 = 4)
		expected := 4 // (1 + 1) * 2
		if result != expected {
			t.Fatalf("expected result %d, got %d", expected, result)
		}
	})

	t.Run("Concurrent Execution", func(t *testing.T) {
		// Arrange
		pipe := NewPipe[int]()

		count := 0

		f := func(ctx context.Context, i int) (int, error) {
			assert.Equal(t, 1, i)
			count++
			return 1 + i, nil
		}

		pipe.Enqueue(f)
		pipe.Enqueue(f)
		pipe.Enqueue(f)

		// Act
		result, err := pipe.Execute(context.Background(), "CONCURRENT", 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Assert

		// Validate that the result is the same passed in the start
		if result != 1 {
			t.Fatalf("expected result to be 1, got %d", result)
		}

		// Validate that the function was executed 3 times. Uses the count variable
		// to count how many times the function was executed
		if count != 3 {
			t.Fatalf("expected count to be 3, got %d", count)
		}
	})

	t.Run("Sequential with First Result Execution", func(t *testing.T) {
		// Arrange
		pipe := NewPipe[int]()
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return 0, assert.AnError
		})
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return i * 3, nil
		})
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return i * 5, nil
		})

		// Act
		result, err := pipe.Execute(context.Background(), SEQUENTIAL_WITH_FIRST_RESULT.String(), 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Assert

		// Validate that the result is 2 * 3 = 6
		// The first function returns an error and is skipped
		// The second function multiplies the input by 3 (2 * 3 = 6)
		// The third function is not executed because the second function succeeded
		expected := 6 // 2 * 3
		if result != expected {
			t.Fatalf("expected result %d, got %d", expected, result)
		}
	})

	t.Run("Skip ObjectNotFound errors in Sequential with First Result Execution", func(t *testing.T) {
		// Arrange
		pipe := NewPipe[int]()
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return 0, ONF
		})
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return i * 4, nil
		})
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return i * 5, nil
		})

		// Act
		result, err := pipe.Execute(context.Background(), SEQUENTIAL_WITH_FIRST_RESULT.String(), 3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Assert

		// Validate that the result is 3 * 4 = 12
		// The first function returns an ObjectNotFound error and is skipped
		// The second function multiplies the input by 4 (3 * 4 = 12)
		// The third function is not executed because the second function succeeded
		expected := 12 // 3 * 4
		if result != expected {
			t.Fatalf("expected result %d, got %d", expected, result)
		}
	})

	t.Run("Should Sequential Blocking Only First Execution", func(t *testing.T) {

		ch := make(chan int)

		// Arrange
		pipe := NewPipe[int]()
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			return i + 10, nil
		})
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			s := i * 5
			assert.Equal(t, 25, s) // Validate that the input is the original input (5)
			time.Sleep(1 * time.Second)
			ch <- s
			return i, nil
		})
		pipe.Enqueue(func(ctx context.Context, i int) (int, error) {
			s := i * 10
			assert.Equal(t, 50, s) // Validate that the input is the original input (5)
			time.Sleep(1 * time.Second)
			ch <- s
			return i, nil
		})

		// Act
		result, err := pipe.Execute(context.Background(), SEQUENTIAL_BLOCKING_ONLY_FIRST.String(), 5)
		assert.NoError(t, err)

		v := 0

		// Wait for the background tasks to complete
		for i := 0; i < 2; i++ {
			v += (<-ch)
		}

		// Assert

		// Validate that the result is 5 + 10 = 15
		// The first function adds 10 to the input (5 + 10 = 15)
		expected := 15 // 5 + 10
		assert.Equal(t, expected, result)

		// Validate that the background tasks were executed
		// The second function multiplies the input by 5 (5 * 5 = 25)
		// The third function multiplies the input by 10 (5 * 10 = 50)
		// The total is 25 + 50 = 75
		assert.Equal(t, 75, v)
	})

}
