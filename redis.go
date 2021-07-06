package simpleq

import (
	"fmt"
	"github.com/go-redis/redis"
)

// NewRedisQueueDriver initializes and returns a pointer to a new redis driver instance
func NewRedisQueueDriver(r redis.Cmdable) *RedisQueueDriver {
	return &RedisQueueDriver{r}
}

// RedisQueueDriver is a queue driver implementation for queue driver
type RedisQueueDriver struct {
	r redis.Cmdable
}

// Write writes to active queue to be executed immediately
func (rqd *RedisQueueDriver) Write(queue string, d []byte) error {
	return rqd.r.SAdd(fmt.Sprintf("%s:active", queue), d).Err()
}

// Read from queue
func (rqd *RedisQueueDriver) Read(queue string) ([]byte, error) {
	return rqd.r.SPop(fmt.Sprintf("%s:active", queue)).Bytes()
}

// SetProcessed increments processed amount
func (rqd *RedisQueueDriver) SetProcessed(queue string) error {
	return rqd.r.Incr(fmt.Sprintf("%s:processed", queue)).Err()
}

// SetFailed increments fail data
func (rqd *RedisQueueDriver) SetFailed(queue string, taskID string) error {
	return rqd.r.SAdd(fmt.Sprintf("%s:failed", queue), taskID).Err()
}

// Register registers a new queue (should not be additive)
func (rqd *RedisQueueDriver) Register(queue string) error {
	return rqd.r.SAdd(fmt.Sprintf("%s:queue-list", queuePrefix), queue).Err()
}

// GetStats returns available queue statistics
func (rqd *RedisQueueDriver) GetStats() (*Stats, error) {
	queues := rqd.r.SMembers(fmt.Sprintf("%s:queue-list", queuePrefix)).Val()

	var stats = Stats{}

	for _, q := range queues {
		proc, _ := rqd.r.Get(fmt.Sprintf("%s:%s:processed", queuePrefix, q)).Int64()
		failed := rqd.r.SMembers(fmt.Sprintf("%s:%s:failed", queuePrefix, q)).Val()

		stats[q] = Stat{
			Processed: proc,
			Failed:    len(failed),
			FailedIDs: failed,
		}
	}

	return &stats, nil
}
