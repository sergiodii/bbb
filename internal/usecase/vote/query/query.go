package query

import (
	"context"
	"fmt"

	usecaseVote "globo_test/internal/usecase/vote"
)

type queryVote struct {
	pipes []pipesExecution
}

type pipesExecution struct {
	pipe          usecaseVote.Pipe[QueryDTO]
	executionType string
	handlerEnum   usecaseVote.HandlerFuncEnum
}

// getPipe retrieves the pipe execution configuration for the given handler function enum.
func (q *queryVote) getPipe(h usecaseVote.HandlerFuncEnum) (pipesExecution, error) {
	for _, p := range q.pipes {
		if p.handlerEnum == h {
			return p, nil
		}
	}
	return pipesExecution{}, fmt.Errorf("pipe not found for handler: %s", h)
}

// QueryVote defines the interface for vote queries.
func (q *queryVote) GetTotalVotes(ctx context.Context, roundID string) (int, error) {
	result, err := q.run(ctx, usecaseVote.HandlerFuncGetTotalVotes, QueryDTO{RoundID: roundID})
	if err != nil || result.Result == nil {
		return 0, err
	}
	return result.Result.(int), nil
}

// GetTotalVotesForParticipant returns a map with the total number of votes for each participant in a given round.
func (q *queryVote) GetTotalVotesForParticipant(ctx context.Context, roundID string) (map[string]int, error) {
	result, err := q.run(ctx, usecaseVote.HandlerFuncGetTotalVotesForParticipant, QueryDTO{RoundID: roundID})
	if err != nil || result.Result == nil {
		return map[string]int{}, err
	}
	return result.Result.(map[string]int), nil
}

// GetTotalVotesForHour returns a map with the total number of votes per hour for a given round.
func (q *queryVote) GetTotalVotesForHour(ctx context.Context, roundID string) (map[string]int, error) {
	result, err := q.run(ctx, usecaseVote.HandlerFuncGetTotalVotesForHour, QueryDTO{RoundID: roundID})
	if err != nil || result.Result == nil {
		return map[string]int{}, err
	}
	return result.Result.(map[string]int), nil
}

// run executes the pipe associated with the given handler function enum.
func (q *queryVote) run(ctx context.Context, handler usecaseVote.HandlerFuncEnum, dto QueryDTO) (QueryDTO, error) {
	p, err := q.getPipe(handler)
	if err != nil {
		return QueryDTO{}, err
	}
	return p.pipe.Execute(ctx, p.executionType, dto)
}

// NewQueryVote creates a new instance of QueryVote with the provided execution pipes.
func NewQueryVote(orderedExecutionPipes map[usecaseVote.HandlerFuncEnum]OrderedExecutionPipeDTO) QueryVoteUseCase {

	pipes := []pipesExecution{}

	for h, executionPipe := range orderedExecutionPipes {
		pipes = append(pipes, pipesExecution{
			pipe:          executionPipe.Pipe,
			executionType: executionPipe.ExecutionType,
			handlerEnum:   h,
		})
	}

	return &queryVote{
		pipes: pipes,
	}
}
