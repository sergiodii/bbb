package api

import (
	"globo_test/cmd/api/route/vote"
	"globo_test/internal/usecase/vote/aggregator"
	"globo_test/pkg/redis"
	"os"

	"github.com/gin-gonic/gin"
)

func commandApiRegister(g *gin.Engine, rootPath string) {
	commandAggregator := aggregator.NewCommandAggregator(
		redis.NewRedisRoundRepository(os.Getenv("REDIS_ADDR")),
		// localsql.NewLocalSqlRoundRepository(), // This is a example, in real case we could have a different repository for command
		// can be added other repositories if needed for example: postgres, mongodb, etc
	)

	vote.NewCommandRoute(commandAggregator, g.Group(rootPath))
}
