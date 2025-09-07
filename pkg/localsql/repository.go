package localsql

import (
	"context"
	"fmt"
	"globo_test/internal/domain/entity"
	"globo_test/internal/domain/repository"
	"sync"
)

var _LocalSqlRoundRepository *LocalSqlRoundRepository
var __LocalSqlRoundRepositoryOnce sync.Once

type LocalSqlRoundRepository struct {
	db []entity.Vote
}

func (lr *LocalSqlRoundRepository) VoteRegister(ctx context.Context, vote entity.Vote) error {
	fmt.Println("Vote registered in Local SQL DB:", vote)
	lr.db = append(lr.db, vote)
	return nil
}

func (lr *LocalSqlRoundRepository) GetTotalVotes(ctx context.Context, roundID string) (int, error) {
	// Implement the logic to get the total votes for a round from the local SQL database
	total := 0
	for _, vote := range lr.db {
		if vote.RoundID == roundID {
			total++
		}
	}

	return total, nil
}

func (lr *LocalSqlRoundRepository) GetTotalForParticipant(ctx context.Context, roundID string) (map[string]int, error) {
	// Implement the logic to get the total votes for each participant in a round from the local SQL database
	total := map[string]int{}
	for _, vote := range lr.db {
		if vote.RoundID == roundID {
			total[vote.ParticipantID]++
		}
	}

	return total, nil
}

func (lr *LocalSqlRoundRepository) GetTotalForHour(ctx context.Context, roundID string) (map[string]int, error) {
	// Implement the logic to get the total votes for each hour in a round from the local SQL database
	total := make(map[string]int)
	for _, vote := range lr.db {
		if vote.RoundID == roundID {
			total[fmt.Sprintf("%d", vote.Timestamp)]++
		}
	}

	return total, nil
}

func NewLocalSqlRoundRepository() repository.RoundRepository {
	__LocalSqlRoundRepositoryOnce.Do(func() {
		_LocalSqlRoundRepository = &LocalSqlRoundRepository{
			db: []entity.Vote{},
		}
	})

	return _LocalSqlRoundRepository
}
