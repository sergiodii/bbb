package command

import (
	"context"
	"globo_test/internal/domain/entity"
)

type CommandVoteUseCase interface {
	CreateVote(ctx context.Context, vote entity.Vote) error
}
