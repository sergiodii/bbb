package vote

import (
	"globo_test/internal/usecase/vote/aggregator"

	"github.com/gin-gonic/gin"
)

func NewQueryRoute(aggregator aggregator.QueryAggregator, g *gin.RouterGroup) {

	queryRoute := newQueryRoute(aggregator.GetAggregatedUseCase())

	g.GET("/:round_id", queryRoute.getTotalVotes())
	g.GET("/:round_id/participant", queryRoute.getTotalVotesForParticipant())
	g.GET("/:round_id/hour", queryRoute.getTotalVotesForHour())
}

func NewCommandRoute(aggregator aggregator.CommandAggregator, g *gin.RouterGroup) {

	commandRoute := newCommandRoute(aggregator.GetAggregatedUseCase())

	g.POST("/:round_id", commandRoute.postCreateVote())
}
