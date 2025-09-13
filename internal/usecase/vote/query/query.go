package query

import (
	"context"

	usecaseVote "github.com/sergiodii/bbb/internal/usecase/vote"
)

type queryVote struct {
	pipeMap map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[QueryDTO]
}

// QueryVote defines the interface for vote queries.
func (q *queryVote) GetTotalVotes(ctx context.Context, roundID string) (int, error) {
	result, err := q.pipeMap[usecaseVote.HandlerFuncGetTotalVotes].Execute(ctx, QueryDTO{RoundID: roundID})
	if err != nil || result.Result == nil {
		return 0, err
	}
	return result.Result.(int), nil
}

// GetTotalVotesForParticipant returns a map with the total number of votes for each participant in a given round.
func (q *queryVote) GetTotalVotesForParticipant(ctx context.Context, roundID string) (map[string]int, error) {
	result, err := q.pipeMap[usecaseVote.HandlerFuncGetTotalVotesForParticipant].Execute(ctx, QueryDTO{RoundID: roundID})
	if err != nil || result.Result == nil {
		return map[string]int{}, err
	}
	return result.Result.(map[string]int), nil
}

// GetTotalVotesForHour returns a map with the total number of votes per hour for a given round.
func (q *queryVote) GetTotalVotesForHour(ctx context.Context, roundID string) (map[string]int, error) {
	result, err := q.pipeMap[usecaseVote.HandlerFuncGetTotalVotesForHour].Execute(ctx, QueryDTO{RoundID: roundID})
	if err != nil || result.Result == nil {
		return map[string]int{}, err
	}
	return result.Result.(map[string]int), nil
}

// NewQueryVote creates a new instance of QueryVote with the provided execution pipes.
func NewQueryVote(pipeMap map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[QueryDTO]) QueryVoteUseCase {
	return &queryVote{
		pipeMap: pipeMap,
	}
}
