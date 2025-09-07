package command

import (
	"context"
	"fmt"
	"globo_test/internal/domain/entity"
	usecaseVote "globo_test/internal/usecase/vote"
)

type commandVote struct {
	pipes []pipesExecution
}

type pipesExecution struct {
	pipe          usecaseVote.Pipe[entity.Vote]
	executionType string
	handlerEnum   usecaseVote.HandlerFuncEnum
}

// getPipe retrieves the pipe execution configuration for the given handler function enum.
func (q *commandVote) getPipe(h usecaseVote.HandlerFuncEnum) (pipesExecution, error) {
	for _, p := range q.pipes {
		if p.handlerEnum == h {
			return p, nil
		}
	}
	return pipesExecution{}, fmt.Errorf("pipe not found for handler: %s", h)
}

// commandVote defines the interface for vote queries.
func (q *commandVote) CreateVote(ctx context.Context, vote entity.Vote) error {
	_, err := q.run(ctx, usecaseVote.HandlerFuncCreateVote, vote)
	return err
}

// run executes the pipe associated with the given handler function enum.
func (q *commandVote) run(ctx context.Context, handler usecaseVote.HandlerFuncEnum, dto entity.Vote) (entity.Vote, error) {
	p, err := q.getPipe(handler)
	if err != nil {
		return entity.Vote{}, err
	}
	return p.pipe.Execute(ctx, p.executionType, dto)
}

// NewCommandVote creates a new instance of commandVote with the provided execution pipes.
func NewCommandVote(orderedExecutionPipes map[usecaseVote.HandlerFuncEnum]OrderedExecutionPipeDTO) CommandVoteUseCase {
	pipes := []pipesExecution{}

	for h, executionPipe := range orderedExecutionPipes {
		pipes = append(pipes, pipesExecution{
			pipe:          executionPipe.Pipe,
			executionType: executionPipe.ExecutionType,
			handlerEnum:   h,
		})
	}

	return &commandVote{
		pipes: pipes,
	}
}
