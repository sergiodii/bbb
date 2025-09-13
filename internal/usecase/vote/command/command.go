package command

import (
	"context"

	"github.com/sergiodii/bbb/internal/domain/entity"
	usecaseVote "github.com/sergiodii/bbb/internal/usecase/vote"
)

type commandVote struct {
	pipeMap map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[entity.Vote]
}

// commandVote defines the interface for vote queries.
func (q *commandVote) CreateVote(ctx context.Context, vote entity.Vote) error {
	_, err := q.pipeMap[usecaseVote.HandlerFuncCreateVote].Execute(ctx, vote)
	return err
}

// NewCommandVote creates a new instance of commandVote with the provided execution pipes.
func NewCommandVote(pipeMap map[usecaseVote.HandlerFuncEnum]usecaseVote.Pipe[entity.Vote]) CommandVoteUseCase {
	return &commandVote{
		pipeMap: pipeMap,
	}
}
