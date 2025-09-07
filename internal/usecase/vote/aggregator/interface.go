package aggregator

import (
	commandVoteUsecase "github.com/sergiodii/bbb/internal/usecase/vote/command"
	queryVoteUsecase "github.com/sergiodii/bbb/internal/usecase/vote/query"
)

type QueryAggregator interface {
	GetAggregatedUseCase() queryVoteUsecase.QueryVoteUseCase
}

type CommandAggregator interface {
	GetAggregatedUseCase() commandVoteUsecase.CommandVoteUseCase
}
