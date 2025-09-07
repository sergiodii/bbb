package vote

import (
	"globo_test/internal/domain/entity"
	commandUsecase "globo_test/internal/usecase/vote/command"
	"time"

	"github.com/gin-gonic/gin"
)

type commandRoute struct {
	uc commandUsecase.CommandVoteUseCase
}

func (q *commandRoute) postCreateVote() func(c *gin.Context) {
	return func(c *gin.Context) {
		roundId := c.Param("round_id")

		var body struct {
			ParticipantID string `json:"participant_id"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		ev := entity.Vote{
			RoundID:       roundId,
			ParticipantID: body.ParticipantID,
			Timestamp:     time.Now().Unix(),
		}

		err := q.uc.CreateVote(c.Request.Context(), ev)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"status": "vote created"})
	}
}

func newCommandRoute(uc commandUsecase.CommandVoteUseCase) *commandRoute {
	return &commandRoute{
		uc: uc,
	}
}
