package vote

import (
	"github.com/gin-gonic/gin"

	queryUsecase "globo_test/internal/usecase/vote/query"
)

type queryRoute struct {
	uc queryUsecase.QueryVoteUseCase
}

func (q *queryRoute) getTotalVotes() func(c *gin.Context) {
	return func(c *gin.Context) {
		pid := c.Param("round_id")

		total, err := q.uc.GetTotalVotes(c.Request.Context(), pid)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"total": total,
		})
	}
}

func (q *queryRoute) getTotalVotesForParticipant() func(c *gin.Context) {
	return func(c *gin.Context) {
		pid := c.Param("round_id")

		totalMap, err := q.uc.GetTotalVotesForParticipant(c.Request.Context(), pid)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, totalMap)
	}
}

func (q *queryRoute) getTotalVotesForHour() func(c *gin.Context) {
	return func(c *gin.Context) {
		pid := c.Param("round_id")

		totalMap, err := q.uc.GetTotalVotesForHour(c.Request.Context(), pid)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, totalMap)
	}
}

func newQueryRoute(uc queryUsecase.QueryVoteUseCase) *queryRoute {
	return &queryRoute{
		uc: uc,
	}
}
