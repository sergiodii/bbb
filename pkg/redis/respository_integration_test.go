package redis_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

	"globo_test/extension/slice"
	"globo_test/internal/domain/entity"
	redisPkg "globo_test/pkg/redis"
)

func TestRedisEstress(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer s.Close()

	repo := redisPkg.NewRedisRoundRepository("localhost:6379")

	entities := slice.MultipliesSlice([]entity.Vote{
		{RoundID: "round1", ParticipantID: "participant1", Timestamp: 1625079600},
		{RoundID: "round1", ParticipantID: "participant2", Timestamp: 1625079600},
		{RoundID: "round1", ParticipantID: "participant1", Timestamp: 1625083200},
		{RoundID: "round1", ParticipantID: "participant3", Timestamp: 1625086800},
	}, 250) // Multiply to create 1000 votes

	counter := struct {
		sync.Mutex
		Count int
	}{}

	wg := sync.WaitGroup{}

	start := time.Now()
	for _, ent := range entities {
		wg.Add(1)
		go func(e entity.Vote) {
			defer wg.Done()
			err := repo.VoteRegister(context.Background(), e)
			if err != nil {
				assert.NoError(t, err)
			}
			counter.Lock()
			counter.Count++
			counter.Unlock()
		}(ent)
	}

	wg.Wait()

	end := time.Now()

	if end.Sub(start) > time.Second*1 {
		t.Fatalf("Vote registration took too long: %v", end.Sub(start))
	}
	t.Logf("Time taken to register 1000 votes: %v", end.Sub(start)*time.Second)

	if counter.Count != 1000 {
		t.Fatalf("Expected %d votes to be registered, got %d", 1000, counter.Count)
	}

	if err != nil {
		t.Fatalf("Failed to register votes: %v", err)
	}

	// total, err := repo.GetTotalVotes(context.Background(), "round1")
	// if err != nil {
	// 	t.Fatalf("Failed to get total votes: %v", err)
	// }
	// if total != 6 {
	// 	t.Fatalf("Expected total votes to be 6, got %d", total)
	// }

	// m, err := repo.GetTotalForParticipant(context.Background(), "round1")
	// if err != nil {
	// 	t.Fatalf("Failed to get total for participant: %v", err)
	// }
	// expected := map[string]int{"participant1": 3, "participant2": 2, "participant3": 1}
	// for k, v := range expected {
	// 	if m[k] != v {
	// 		t.Fatalf("Expected participant %s to have %d votes, got %d", k, v, m[k])
	// 	}
	// }

	// h, err := repo.GetTotalForHour(context.Background(), "round1")
	// if err != nil {
	// 	t.Fatalf("Failed to get total for hour: %v", err)
	// }
	// expectedHours := map[string]int{"14": 2, "15": 2, "16": 1, "17": 1}
	// for k, v := range expectedHours {
	// 	if h[k] != v {
	// 		t.Fatalf("Expected hour %s to have %d votes, got %d", k, v, h[k])
	// 	}
	// }

}
