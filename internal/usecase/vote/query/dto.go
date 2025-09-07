package query

import usecaseVote "globo_test/internal/usecase/vote"

type QueryDTO struct {
	RoundID       string
	ParticipantID string
	Result        interface{}
}

type OrderedExecutionPipeDTO struct {
	ExecutionType string
	Pipe          usecaseVote.Pipe[QueryDTO]
}
