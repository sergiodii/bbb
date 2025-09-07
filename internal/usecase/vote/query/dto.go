package query

import usecaseVote "github.com/sergiodii/bbb/internal/usecase/vote"

type QueryDTO struct {
	RoundID       string
	ParticipantID string
	Result        interface{}
}

type OrderedExecutionPipeDTO struct {
	ExecutionType string
	Pipe          usecaseVote.Pipe[QueryDTO]
}
