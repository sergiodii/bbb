package query_test

import (
	"context"
	"testing"

	"github.com/sergiodii/bbb/internal/usecase/vote/query"

	"github.com/stretchr/testify/assert"

	usecaseVote "github.com/sergiodii/bbb/internal/usecase/vote"
)

type pipeMock struct {
	ExecutedFuncs []func(context.Context, query.QueryDTO) (query.QueryDTO, error)
}

func (p *pipeMock) Enqueue(funcs ...func(context.Context, query.QueryDTO) (query.QueryDTO, error)) {
	p.ExecutedFuncs = append(p.ExecutedFuncs, funcs...)
}

func (p *pipeMock) Execute(ctx context.Context, executionType string, dto query.QueryDTO) (query.QueryDTO, error) {
	var err error
	for _, fn := range p.ExecutedFuncs {
		dto, err = fn(ctx, dto)
		if err != nil {
			return dto, err
		}
	}
	return dto, nil
}

func TestGetVotesFromParticipant(t *testing.T) {

	t.Run("Should return total votes", func(t *testing.T) {
		pm := &pipeMock{}
		pm.Enqueue(func(ctx context.Context, dto query.QueryDTO) (query.QueryDTO, error) {
			dto.Result = 1 * 5
			return dto, nil
		})

		execution := map[usecaseVote.HandlerFuncEnum]query.OrderedExecutionPipeDTO{
			usecaseVote.HandlerFuncGetTotalVotes: {
				ExecutionType: "SEQUENTIAL",
				Pipe:          pm,
			},
		}

		q := query.NewQueryVote(execution)
		result, err := q.GetTotalVotes(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, 5, result)
	})

	t.Run("Should return total votes for participant", func(t *testing.T) {
		pm := &pipeMock{}
		pm.Enqueue(func(ctx context.Context, dto query.QueryDTO) (query.QueryDTO, error) {
			dto.Result = map[string]int{"participant1": 10, "participant2": 5}
			return dto, nil
		})

		execution := map[usecaseVote.HandlerFuncEnum]query.OrderedExecutionPipeDTO{
			usecaseVote.HandlerFuncGetTotalVotesForParticipant: {
				ExecutionType: "SEQUENTIAL",
				Pipe:          pm,
			},
		}

		q := query.NewQueryVote(execution)
		result, err := q.GetTotalVotesForParticipant(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, map[string]int{"participant1": 10, "participant2": 5}, result)
	})

	t.Run("Should return total votes for participant summed", func(t *testing.T) {
		pm := &pipeMock{}
		pm.Enqueue(func(ctx context.Context, dto query.QueryDTO) (query.QueryDTO, error) {
			dto.Result = map[string]int{"participant1": 10, "participant2": 5}
			return dto, nil
		})
		pm.Enqueue(func(ctx context.Context, dto query.QueryDTO) (query.QueryDTO, error) {
			currentResult := dto.Result.(map[string]int)
			currentResult["participant1"] += 15
			currentResult["participant3"] = 7
			dto.Result = currentResult
			return dto, nil
		})

		execution := map[usecaseVote.HandlerFuncEnum]query.OrderedExecutionPipeDTO{
			usecaseVote.HandlerFuncGetTotalVotesForParticipant: {
				ExecutionType: "SEQUENTIAL",
				Pipe:          pm,
			},
		}

		q := query.NewQueryVote(execution)
		result, err := q.GetTotalVotesForParticipant(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, map[string]int{"participant1": 25, "participant2": 5, "participant3": 7}, result)
	})

	t.Run("Should handle pipe execution error", func(t *testing.T) {

		pm := &pipeMock{}
		pm.Enqueue(func(ctx context.Context, dto query.QueryDTO) (query.QueryDTO, error) {
			return dto, assert.AnError
		})

		execution := map[usecaseVote.HandlerFuncEnum]query.OrderedExecutionPipeDTO{
			usecaseVote.HandlerFuncGetTotalVotes: {
				ExecutionType: "SEQUENTIAL",
				Pipe:          pm,
			},
		}

		q := query.NewQueryVote(execution)
		_, err := q.GetTotalVotes(context.Background(), "1")
		assert.Error(t, err)
	})

	t.Run("Should handle missing pipe", func(t *testing.T) {

		execution := map[usecaseVote.HandlerFuncEnum]query.OrderedExecutionPipeDTO{}

		q := query.NewQueryVote(execution)
		_, err := q.GetTotalVotes(context.Background(), "1")
		assert.Error(t, err)
	})
}
