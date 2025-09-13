package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sergiodii/bbb/extension/slice"
	"github.com/sergiodii/bbb/extension/text"
	"github.com/sergiodii/bbb/internal/domain/entity"
	"github.com/sergiodii/bbb/internal/domain/repository"

	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
)

type totalResult struct {
	Total int
	m     sync.Mutex
}

func (t *totalResult) Add(v int) {
	t.m.Lock()
	defer t.m.Unlock()
	t.Total += v
}

type RedisRoundRepository struct {
	Client *redis.Client
}

// VoteRegister registers a vote in Redis by incrementing the count for the participant and the hour.
// It uses errgroup to perform both increments concurrently.
// When occurs an error, a rollback policy should be implemented to ensure data consistency. A method example to decrement the counts could be VoteRegisterRollback(ctx context.Context, vote entity.Vote) error.
// However, for simplicity, this example does not include rollback logic.
func (r *RedisRoundRepository) VoteRegister(ctx context.Context, vote entity.Vote) error {

	var eg errgroup.Group

	eg.Go(func() error {
		key := fmt.Sprintf("round:%s:participant:%s", vote.RoundID, vote.ParticipantID)
		return r.Client.Incr(ctx, key).Err()
	})

	eg.Go(func() error {
		hourKey := fmt.Sprintf("round:%s:hour:%d", vote.RoundID, vote.Timestamp/3600)
		return r.Client.Incr(ctx, hourKey).Err()
	})

	eg.Go(func() error {
		totalKey := fmt.Sprintf("round:%s:total", vote.RoundID)
		return r.Client.Incr(ctx, totalKey).Err()
	})

	if err := eg.Wait(); err != nil {

		// rollback logic should be implemented here to ensure data consistency
		// for simplicity, it's omitted in this example
		return err
	}

	return nil
}

func (r *RedisRoundRepository) GetTotalVotes(ctx context.Context, roundID string) (int, error) {
	key := fmt.Sprintf("round:%s:total", roundID)
	val, err := r.Client.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *RedisRoundRepository) GetTotalForParticipant(ctx context.Context, roundID string) (map[string]int, error) {
	pattern := fmt.Sprintf("round:%s:participant:*", roundID)
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	result := &sync.Map{}
	var eg errgroup.Group
	for _, chunk := range slice.TransformSliceToMultipleSlices(keys, 10) {
		eg.Go(func() error {
			for _, key := range chunk {
				v, err := r.Client.Get(ctx, key).Int()
				if err != nil {
					return err
				}
				// key format: round:<id>:participant:<pid>
				pid := text.GetPartOfString(key, `round:[^:]+:participant:(.+)`)
				result.Store(pid, v)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return convertSyncMapToMapStringInt(result), nil
}

func (r *RedisRoundRepository) GetTotalForHour(ctx context.Context, roundID string) (map[string]int, error) {
	pattern := fmt.Sprintf("round:%s:hour:*", roundID)
	keys, err := r.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	result := make(map[string]int)
	for _, key := range keys {
		v, err := r.Client.Get(ctx, key).Int()
		if err == nil {
			// key format: round:<id>:hour:<hour>
			hour := text.GetPartOfString(key, `round:[^:]+:hour:(.+)`)
			result[hour] = v
		}
	}
	return result, nil
}

func convertSyncMapToMapStringInt(sm *sync.Map) map[string]int {
	result := make(map[string]int)
	sm.Range(func(key, value any) bool {
		k, ok1 := key.(string)
		v, ok2 := value.(int)
		if ok1 && ok2 {
			result[k] = v
		}
		return true
	})
	return result
}

func NewRedisRoundRepository(addr string) repository.RoundRepository {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     100,             // Aumentar pool de conexões para alta concorrência
		MinIdleConns: 10,              // Manter conexões ativas para reduzir latência
		MaxRetries:   3,               // Retry em caso de timeout
		DialTimeout:  5 * time.Second, // Timeout para conectar
		ReadTimeout:  3 * time.Second, // Timeout para leitura
		WriteTimeout: 3 * time.Second, // Timeout para escrita
	})
	return &RedisRoundRepository{Client: client}
}
