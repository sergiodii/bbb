package command

import (
	"context"

	"github.com/sergiodii/bbb/internal/domain/entity"
)

type CommandVoteUseCase interface {
	CreateVote(ctx context.Context, vote entity.Vote) error
}
