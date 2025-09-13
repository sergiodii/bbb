package command

import (
	"context"
	"testing"

	"github.com/sergiodii/bbb/internal/domain/entity"
	usecaseVote "github.com/sergiodii/bbb/internal/usecase/vote"
	"github.com/sergiodii/bbb/internal/usecase/vote/query/mock"

	"github.com/stretchr/testify/assert"
)

func TestNewCommandVote(t *testing.T) {

	t.Run("Should create a new CommandVote instance", func(t *testing.T) {

		// Arrange
		pipe := mock.NewPipeMock[entity.Vote]()

		pipeMap := map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[entity.Vote]{
			usecaseVote.HandlerFuncCreateVote: pipe,
		}

		// Act
		commandVote := NewCommandVote(pipeMap)

		// Assert
		if commandVote == nil {
			t.Fatal("Expected non-nil CommandVote")
		}
	})

	t.Run("Should handle empty orderedExecutionPipes", func(t *testing.T) {

		// Arrange
		pipeMap := map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[entity.Vote]{}

		// Act
		commandVote := NewCommandVote(pipeMap)

		// Assert
		if commandVote == nil {
			t.Fatal("Expected non-nil CommandVote even with empty pipes")
		}
	})

	t.Run("Should execute CreateVote without error", func(t *testing.T) {

		// Arrange
		pipe := mock.NewPipeMock[entity.Vote]()

		e := entity.Vote{RoundID: "round1", ParticipantID: "participant1", Timestamp: 1234567890}

		pipe.On("Execute", context.Background(), e).Return(e, nil)

		pipeMap := map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[entity.Vote]{
			usecaseVote.HandlerFuncCreateVote: pipe,
		}

		commandVote := NewCommandVote(pipeMap)

		// Act
		err := commandVote.CreateVote(context.Background(), entity.Vote{RoundID: "round1", ParticipantID: "participant1", Timestamp: e.Timestamp})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "round1", "round1")
		assert.Equal(t, "participant1", "participant1")
		assert.Equal(t, 1234567890, 1234567890)
	})
}
