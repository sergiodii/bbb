package repository

import (
	"context"

	"globo_test/internal/domain/entity"
)

type RoundRepository interface {
	VoteRegister(ctx context.Context, vote entity.Vote) error
	GetTotalVotes(ctx context.Context, roundID string) (int, error)
	GetTotalForParticipant(ctx context.Context, roundID string) (map[string]int, error)
	GetTotalForHour(ctx context.Context, roundID string) (map[string]int, error)
}
