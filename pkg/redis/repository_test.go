package redis

import (
	"context"
	"testing"

	"github.com/sergiodii/bbb/internal/domain/entity"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestVoteRegister_Success(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer s.Close()

	repo := NewRedisRoundRepository(s.Addr())

	err = repo.VoteRegister(context.Background(), entity.Vote{
		RoundID:       "round1",
		ParticipantID: "participant1",
		Timestamp:     1625079600, // Example timestamp
	})
	assert.NoError(t, err)

	total, err := repo.GetTotalVotes(context.Background(), "round1")
	assert.Equal(t, 1, total)
	assert.NoError(t, err)

	m, err := repo.GetTotalForParticipant(context.Background(), "round1")
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"participant1": 1}, m)
}
