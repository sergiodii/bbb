package command_test

import (
	"context"
	"globo_test/internal/domain/entity"
	"globo_test/internal/usecase/vote/command"
	"testing"

	"github.com/stretchr/testify/assert"

	usecaseVote "globo_test/internal/usecase/vote"
)

type pipeMock struct {
	ExecutedFuncs []func(context.Context, entity.Vote) (entity.Vote, error)
}

func (p *pipeMock) Enqueue(funcs ...func(context.Context, entity.Vote) (entity.Vote, error)) {
	p.ExecutedFuncs = append(p.ExecutedFuncs, funcs...)
}

func (p *pipeMock) Execute(ctx context.Context, executionType string, dto entity.Vote) (entity.Vote, error) {
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

	t.Run("Should set vote", func(t *testing.T) {
		pm := &pipeMock{}
		pm.Enqueue(func(ctx context.Context, dto entity.Vote) (entity.Vote, error) {
			return dto, nil
		})

		execution := map[usecaseVote.HandlerFuncEnum]command.OrderedExecutionPipeDTO{
			usecaseVote.HandlerFuncCreateVote: {
				ExecutionType: "SEQUENTIAL",
				Pipe:          pm,
			},
		}

		q := command.NewCommandVote(execution)
		err := q.CreateVote(context.Background(), entity.Vote{RoundID: "round1", ParticipantID: "participant1", Timestamp: 1234567890})
		assert.NoError(t, err)
	})

	t.Run("Should handle missing pipe", func(t *testing.T) {

		execution := map[usecaseVote.HandlerFuncEnum]command.OrderedExecutionPipeDTO{}

		q := command.NewCommandVote(execution)
		err := q.CreateVote(context.Background(), entity.Vote{RoundID: "round1", ParticipantID: "participant1", Timestamp: 1234567890})
		assert.Error(t, err)
	})
}
