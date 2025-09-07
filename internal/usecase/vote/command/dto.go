package command

import (
	"globo_test/internal/domain/entity"
	usecaseVote "globo_test/internal/usecase/vote"
)

type OrderedExecutionPipeDTO struct {
	ExecutionType string
	Pipe          usecaseVote.Pipe[entity.Vote]
}
