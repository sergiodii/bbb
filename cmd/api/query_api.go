package api

import (
	"os"

	"github.com/sergiodii/bbb/cmd/api/route/vote"
	"github.com/sergiodii/bbb/internal/usecase/vote/aggregator"
	"github.com/sergiodii/bbb/pkg/redis"

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
