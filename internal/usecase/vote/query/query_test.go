package query

import (
	"context"
	"testing"

	"github.com/sergiodii/bbb/internal/usecase/vote/query/mock"

	usecaseVote "github.com/sergiodii/bbb/internal/usecase/vote"
)

func TestNewQueryVote(t *testing.T) {

	t.Run("Should create a new QueryVote instance", func(t *testing.T) {

		// Arrange
		pipe := mock.NewPipeMock[QueryDTO]()

		pipeMap := map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[QueryDTO]{
			usecaseVote.HandlerFuncGetTotalVotes: pipe,
		}

		// Act
		queryVote := NewQueryVote(pipeMap)

		// Assert
		if queryVote == nil {
			t.Fatal("Expected non-nil QueryVote")
		}
	})

	t.Run("Should handle empty orderedExecutionPipes", func(t *testing.T) {

		// Arrange
		orderedExecutionPipes := map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[QueryDTO]{}

		// Act
		queryVote := NewQueryVote(orderedExecutionPipes)

		// Assert
		if queryVote == nil {
			t.Fatal("Expected non-nil QueryVote even with empty pipes")
		}
	})

	t.Run("Should execute GetTotalVotes without error", func(t *testing.T) {

		// Arrange
		pipe := mock.NewPipeMock[QueryDTO]()

		pipe.On("Execute", context.Background(), QueryDTO{RoundID: "round1"}).Return(QueryDTO{Result: 42}, nil)

		orderedExecutionPipes := map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[QueryDTO]{
			usecaseVote.HandlerFuncGetTotalVotes: pipe,
		}

		queryVote := NewQueryVote(orderedExecutionPipes)

		// Act
		totalVotes, err := queryVote.GetTotalVotes(context.Background(), "round1")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if totalVotes != 42 {
			t.Fatalf("Expected totalVotes to be 42, got %d", totalVotes)
		}
	})

	t.Run("Should execute GetTotalVotesForParticipant without error", func(t *testing.T) {

		// Arrange
		pipe := mock.NewPipeMock[QueryDTO]()

		expectedResult := map[string]int{"participant1": 10, "participant2": 20}
		pipe.On("Execute", context.Background(), QueryDTO{RoundID: "round1"}).Return(QueryDTO{Result: expectedResult}, nil)

		orderedExecutionPipes := map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[QueryDTO]{
			usecaseVote.HandlerFuncGetTotalVotesForParticipant: pipe,
		}

		queryVote := NewQueryVote(orderedExecutionPipes)

		// Act
		result, err := queryVote.GetTotalVotesForParticipant(context.Background(), "round1")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(result) != len(expectedResult) {
			t.Fatalf("Expected result length %d, got %d", len(expectedResult), len(result))
		}
		for k, v := range expectedResult {
			if result[k] != v {
				t.Fatalf("Expected result[%s] to be %d, got %d", k, v, result[k])
			}
		}
	})

	t.Run("Should execute GetTotalVotesForHour without error", func(t *testing.T) {
		// Arrange
		pipe := mock.NewPipeMock[QueryDTO]()

		expectedResult := map[string]int{"10:00": 5, "11:00": 15}
		pipe.On("Execute", context.Background(), QueryDTO{RoundID: "round1"}).Return(QueryDTO{Result: expectedResult}, nil)

		orderedExecutionPipes := map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[QueryDTO]{
			usecaseVote.HandlerFuncGetTotalVotesForHour: pipe,
		}

		queryVote := NewQueryVote(orderedExecutionPipes)

		// Act
		result, err := queryVote.GetTotalVotesForHour(context.Background(), "round1")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(result) != len(expectedResult) {
			t.Fatalf("Expected result length %d, got %d", len(expectedResult), len(result))
		}
		for k, v := range expectedResult {
			if result[k] != v {
				t.Fatalf("Expected result[%s] to be %d, got %d", k, v, result[k])
			}
		}
	})
}
