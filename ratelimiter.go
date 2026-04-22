package shortly

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type RateLimiterStore struct {
	limiters map[string]*IPLimiter
	mu       sync.Mutex
}

type IPLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RedisRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

type RateLimiterAllow interface {
	Allow(string) (bool, error)
}

func RedisRateLimiterStore() *RedisRateLimiter {
	redisClient := GetRedisConnection()
	return &RedisRateLimiter{
		client: redisClient,
		limit:  6,
		window: time.Second * 10,
	}
}
func (rr *RedisRateLimiter) Allow(ip string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit_%s", ip)
	count, err := rr.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if count == 1 {
		rr.client.Expire(ctx, key, rr.window)
	}
	return count <= int64(rr.limit), nil

}

func NewRateLimiterStore() *RateLimiterStore {
	return &RateLimiterStore{
		limiters: make(map[string]*IPLimiter),
		mu:       sync.Mutex{},
	}
}

func (rl *RateLimiterStore) getIpLimiter(ip string) *IPLimiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	result, ok := rl.limiters[ip]
	if !ok {
		result = &IPLimiter{
			limiter:  rate.NewLimiter(1, 2),
			lastSeen: time.Now(),
		}
		rl.limiters[ip] = result
	} else {
		result.lastSeen = time.Now()
	}

	return result

}

func (rl *RateLimiterStore) RateLimiterCleanUp() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		for ip, l := range rl.limiters {
			if time.Since(l.lastSeen) > 10*time.Minute {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}
