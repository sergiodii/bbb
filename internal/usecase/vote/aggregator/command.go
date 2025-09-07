package aggregator

import (
	"context"
	"sync"

	"globo_test/extension/pipe"
	"globo_test/internal/domain/entity"
	"globo_test/internal/domain/repository"
	voteUsecase "globo_test/internal/usecase/vote"
	commandVoteUsecase "globo_test/internal/usecase/vote/command"
)

var commandAggregated *commandAggregator

var commandAggregatedOnce sync.Once

type commandAggregator struct {
	repositories []repository.RoundRepository
}

func (a *commandAggregator) aggregateVoteRegisterHandler() pipe.Pipe[entity.Vote] {
	p := pipe.NewPipe[entity.Vote]()
	for _, exec := range a.repositories {
		p.Enqueue(func(ctx context.Context, dto entity.Vote) (entity.Vote, error) {
			err := exec.VoteRegister(ctx, dto)
			if err != nil {
				return dto, err
			}
			return dto, nil
		})
	}
	return p
}

func (a *commandAggregator) GetAggregatedUseCase() commandVoteUsecase.CommandVoteUseCase {

	executionMap := map[voteUsecase.HandlerFuncEnum]commandVoteUsecase.OrderedExecutionPipeDTO{
		voteUsecase.HandlerFuncCreateVote: {
			ExecutionType: pipe.SEQUENTIAL_BLOCKING_ONLY_FIRST.String(),
			Pipe:          a.aggregateVoteRegisterHandler(),
		},
	}
	return commandVoteUsecase.NewCommandVote(executionMap)
}

func NewCommandAggregator(repos ...repository.RoundRepository) CommandAggregator {

	commandAggregatedOnce.Do(func() {
		commandAggregated = &commandAggregator{
			repositories: repos,
		}
	})

	return commandAggregated
}
