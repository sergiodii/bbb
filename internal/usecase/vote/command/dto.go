package command

import (
	"github.com/sergiodii/bbb/internal/domain/entity"
	usecaseVote "github.com/sergiodii/bbb/internal/usecase/vote"
)

type OrderedExecutionPipeDTO struct {
	ExecutionType string
	Pipe          usecaseVote.Pipe[entity.Vote]
}
