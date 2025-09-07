package aggregator

import (
	"context"
	"globo_test/extension/pipe"
	"globo_test/internal/domain/repository"
	voteUsecase "globo_test/internal/usecase/vote"
	queryVoteUsecase "globo_test/internal/usecase/vote/query"

	"sync"
)

var queryAggregated *queryAggregator

var queryAggregatedOnce sync.Once

type queryAggregator struct {
	repositories []repository.RoundRepository
}

func (a *queryAggregator) aggregateTotalVotesHandler() pipe.Pipe[queryVoteUsecase.QueryDTO] {
	p := pipe.NewPipe[queryVoteUsecase.QueryDTO]()
	for _, exec := range a.repositories {
		p.Enqueue(func(ctx context.Context, dto queryVoteUsecase.QueryDTO) (queryVoteUsecase.QueryDTO, error) {
			total, err := exec.GetTotalVotes(ctx, dto.RoundID)
			if err != nil {
				return dto, err
			}

			if total == 0 {
				// If no votes found, return ObjectNotFound error to let the pipe continue
				return dto, pipe.ONF
			}

			dto.Result = total
			return dto, nil
		})
	}
	return p
}

func (a *queryAggregator) aggregateTotalVotesForParticipantHandler() pipe.Pipe[queryVoteUsecase.QueryDTO] {
	p := pipe.NewPipe[queryVoteUsecase.QueryDTO]()
	for _, exec := range a.repositories {
		p.Enqueue(func(ctx context.Context, dto queryVoteUsecase.QueryDTO) (queryVoteUsecase.QueryDTO, error) {
			totalMap, err := exec.GetTotalForParticipant(ctx, dto.RoundID)

			if len(totalMap) == 0 {

				// If no votes found, return ObjectNotFound error to let the pipe continue
				return dto, pipe.ONF
			}
			if err != nil {
				return dto, err
			}
			dto.Result = totalMap
			return dto, nil
		})
	}
	return p
}

func (a *queryAggregator) aggregateTotalVotesForHourHandler() pipe.Pipe[queryVoteUsecase.QueryDTO] {
	p := pipe.NewPipe[queryVoteUsecase.QueryDTO]()
	for _, exec := range a.repositories {
		p.Enqueue(func(ctx context.Context, dto queryVoteUsecase.QueryDTO) (queryVoteUsecase.QueryDTO, error) {
			totalMap, err := exec.GetTotalForHour(ctx, dto.RoundID)

			if len(totalMap) == 0 {
				// If no votes found, return ObjectNotFound error to let the pipe continue
				return dto, pipe.ONF
			}

			if err != nil {
				return dto, err
			}

			dto.Result = totalMap
			return dto, nil
		})
	}
	return p
}

func (a *queryAggregator) GetAggregatedUseCase() queryVoteUsecase.QueryVoteUseCase {

	executionMap := map[voteUsecase.HandlerFuncEnum]queryVoteUsecase.OrderedExecutionPipeDTO{
		voteUsecase.HandlerFuncGetTotalVotes: {
			ExecutionType: pipe.SEQUENTIAL_WITH_FIRST_RESULT.String(),
			Pipe:          a.aggregateTotalVotesHandler(),
		},
		voteUsecase.HandlerFuncGetTotalVotesForParticipant: {
			ExecutionType: pipe.SEQUENTIAL_WITH_FIRST_RESULT.String(),
			Pipe:          a.aggregateTotalVotesForParticipantHandler(),
		},
		voteUsecase.HandlerFuncGetTotalVotesForHour: {
			ExecutionType: pipe.SEQUENTIAL_WITH_FIRST_RESULT.String(),
			Pipe:          a.aggregateTotalVotesForHourHandler(),
		},
	}
	return queryVoteUsecase.NewQueryVote(executionMap)
}

func NewQueryAggregator(repos ...repository.RoundRepository) QueryAggregator {

	queryAggregatedOnce.Do(func() {
		queryAggregated = &queryAggregator{
			repositories: repos,
		}
	})

	return queryAggregated
}
