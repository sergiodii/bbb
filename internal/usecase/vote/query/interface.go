package query

import "context"

type QueryVoteUseCase interface {

	// Returns the total number of votes for a given round.
	GetTotalVotes(ctx context.Context, roundID string) (int, error)

	// Returns a map with the total number of votes for each participant in a given round.
	GetTotalVotesForParticipant(ctx context.Context, roundID string) (map[string]int, error)

	// Returns a map with the total number of votes per hour for a given round.
	GetTotalVotesForHour(ctx context.Context, roundID string) (map[string]int, error)
}
