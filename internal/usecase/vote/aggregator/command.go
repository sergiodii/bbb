package aggregator

import (
	"context"
	"sync"

	"github.com/sergiodii/bbb/extension/pipe"
	"github.com/sergiodii/bbb/internal/domain/entity"
	"github.com/sergiodii/bbb/internal/domain/repository"
	voteUsecase "github.com/sergiodii/bbb/internal/usecase/vote"
	commandVoteUsecase "github.com/sergiodii/bbb/internal/usecase/vote/command"
)

var commandAggregated *commandAggregator

var commandAggregatedOnce sync.Once

type commandAggregator struct {
	repositories []repository.RoundRepository
}

func (a *commandAggregator) aggregateVoteRegisterHandler() pipe.Pipe[entity.Vote] {
	p := pipe.NewPipe[entity.Vote](pipe.SEQUENTIAL_BLOCKING_ONLY_FIRST)
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

	executionMap := map[voteUsecase.HandlerFuncEnum]voteUsecase.Pipe[entity.Vote]{
		voteUsecase.HandlerFuncCreateVote: a.aggregateVoteRegisterHandler(),
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
