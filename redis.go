package simpleq

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// NewRedisQueueDriver initializes and returns a pointer to a new redis driver instance
func NewRedisQueueDriver(r RedisInterface) *RedisQueueDriver {
	return &RedisQueueDriver{r}
}

// RedisInterface is an interface for needed methods of the redis client [go get github.com/go-redis/redis]
type RedisInterface interface {
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SPop(ctx context.Context, key string) *redis.StringCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	//Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
}

// RedisQueueDriver is a queue driver implementation for queue driver
type RedisQueueDriver struct {
	r RedisInterface
}

// Write writes to active queue to be executed immediately
func (rqd *RedisQueueDriver) Write(queue string, d []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return rqd.r.SAdd(ctx, fmt.Sprintf("%s:active", queue), d).Err()
}

// Read from queue
func (rqd *RedisQueueDriver) Read(queue string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return rqd.r.SPop(ctx, fmt.Sprintf("%s:active", queue)).Bytes()
}

// SetProcessed increments processed amount
func (rqd *RedisQueueDriver) SetProcessed(queue string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return rqd.r.Incr(ctx, fmt.Sprintf("%s:processed", queue)).Err()
}

// SetFailed increments fail data
func (rqd *RedisQueueDriver) SetFailed(queue string, taskID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return rqd.r.SAdd(ctx, fmt.Sprintf("%s:failed", queue), taskID).Err()
}

// Register registers a new queue (should not be additive)
func (rqd *RedisQueueDriver) Register(queue string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return rqd.r.SAdd(ctx, fmt.Sprintf("%s:queue-list", queuePrefix), queue).Err()
}

// GetStats returns available queue statistics
func (rqd *RedisQueueDriver) GetStats() (*Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	queues := rqd.r.SMembers(ctx, fmt.Sprintf("%s:queue-list", queuePrefix)).Val()

	var stats = Stats{}

	for _, q := range queues {
		proc, _ := rqd.r.Get(ctx, fmt.Sprintf("%s:%s:processed", queuePrefix, q)).Int64()
		failed := rqd.r.SMembers(ctx, fmt.Sprintf("%s:%s:failed", queuePrefix, q)).Val()

		stats[q] = Stat{
			Processed: proc,
			Failed:    len(failed),
			FailedIDs: failed,
		}
	}

	return &stats, nil
}
