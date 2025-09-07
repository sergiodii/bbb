package aggregator

import (
	commandVoteUsecase "globo_test/internal/usecase/vote/command"
	queryVoteUsecase "globo_test/internal/usecase/vote/query"
)

type QueryAggregator interface {
	GetAggregatedUseCase() queryVoteUsecase.QueryVoteUseCase
}

type CommandAggregator interface {
	GetAggregatedUseCase() commandVoteUsecase.CommandVoteUseCase
}
