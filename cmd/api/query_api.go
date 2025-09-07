package api

import (
	"globo_test/cmd/api/route/vote"
	"globo_test/internal/usecase/vote/aggregator"
	"globo_test/pkg/redis"
	"os"

	"github.com/gin-gonic/gin"
)

func queryApiRegister(g *gin.Engine, rootPath string) {
	queryAggregator := aggregator.NewQueryAggregator(
		redis.NewRedisRoundRepository(os.Getenv("REDIS_ADDR")),
		// localsql.NewLocalSqlRoundRepository(), // This is a example, in real case we could have a different repository for command
		// can be added other repositories if needed for example: postgres, mongodb, etc
	)

	vote.NewQueryRoute(queryAggregator, g.Group(rootPath))
}
