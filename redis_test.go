package simpleq

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestRedisQueueDriver_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := NewMockRedisInterface(ctrl)

	qName := "test-queue"

	t.Run("it_should_return_error_when_it_fails_to_write_data", func(t *testing.T) {
		expectErr := fmt.Errorf("failed to add data")

		r.
			EXPECT().
			SAdd(gomock.Any(), "test-queue:active", gomock.Any()).
			DoAndReturn(func(_ interface{}, _ string, _ interface{}) *redis.IntCmd {
				return redis.NewIntResult(0, fmt.Errorf("failed to add data"))
			}).
			Times(1)

		d := NewRedisQueueDriver(r)

		if err := d.Write(qName, []byte("test")); err == nil || err.Error() != expectErr.Error() {
			t.Errorf("Expected Write() to return %v, got %v", expectErr, err)
		}
	})

	t.Run("it_should_write_data", func(t *testing.T) {
		r.
			EXPECT().
			SAdd(gomock.Any(), "test-queue:active", gomock.Any()).
			DoAndReturn(func(_ interface{}, _ string, _ interface{}) *redis.IntCmd {
				return redis.NewIntResult(1, nil)
			}).
			Times(1)

		d := NewRedisQueueDriver(r)

		if err := d.Write(qName, []byte("test")); err != nil {
			t.Errorf("Expected Write() to write data, got error %v", err)
		}
	})
}

func TestRedisQueueDriver_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := NewMockRedisInterface(ctrl)

	qName := "test-queue"

	t.Run("it_should_return_error_when_it_fails_to_read_data", func(t *testing.T) {
		expect := fmt.Errorf("failed to fetch")

		r.
			EXPECT().
			SPop(gomock.Any(), "test-queue:active").
			DoAndReturn(func(_ interface{}, _ string) *redis.StringCmd {
				return redis.NewStringResult("", fmt.Errorf("failed to fetch"))
			}).
			Times(1)

		d := NewRedisQueueDriver(r)

		if _, err := d.Read(qName); err == nil || err.Error() != expect.Error() {
			t.Errorf("Expected Read() to return fetch error, got %v", err)
		}
	})
}

func TestRedisQueueDriver_SetProcessed(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := NewMockRedisInterface(ctrl)

	qName := "test-queue"

	t.Run("it_should_return_error_when_it_fails_to_set_processed", func(t *testing.T) {
		expect := fmt.Errorf("failed to increment")

		r.
			EXPECT().
			Incr(gomock.Any(), "test-queue:processed").
			DoAndReturn(func(_ interface{}, _ string) *redis.IntCmd {
				return redis.NewIntResult(0, fmt.Errorf("failed to increment"))
			}).
			Times(1)

		d := NewRedisQueueDriver(r)

		if err := d.SetProcessed(qName); err == nil || err.Error() != expect.Error() {
			t.Errorf("Expected SetProcessed() to return %v, got %v", expect, err)
		}
	})

	t.Run("it_should_increment_processed_count", func(t *testing.T) {
		r.
			EXPECT().
			Incr(gomock.Any(), "test-queue:processed").
			DoAndReturn(func(_ interface{}, _ string) *redis.IntCmd {
				return redis.NewIntResult(1, nil)
			}).
			Times(1)

		d := NewRedisQueueDriver(r)

		if err := d.SetProcessed(qName); err != nil {
			t.Errorf("Expected SetProcessed() to increment data, got %v", err)
		}
	})
}

func TestRedisQueueDriver_SetFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := NewMockRedisInterface(ctrl)

	qName := "test-queue"
	taskId := uuid.New().String()

	t.Run("it_should_return_error_when_it_fails_to_set_failed", func(t *testing.T) {
		expect := fmt.Errorf("failed to write")

		r.
			EXPECT().
			SAdd(gomock.Any(), "test-queue:failed", taskId).
			DoAndReturn(func(_ interface{}, _, _ string) *redis.IntCmd {
				return redis.NewIntResult(0, fmt.Errorf("failed to write"))
			}).
			Times(1)

		if err := NewRedisQueueDriver(r).SetFailed(qName, taskId); err == nil || err.Error() != expect.Error() {
			t.Errorf("Expected SetFailed() to return %v, got %v", expect, err)
		}
	})

	t.Run("it_should_set_failed", func(t *testing.T) {
		r.
			EXPECT().
			SAdd(gomock.Any(), "test-queue:failed", taskId).
			DoAndReturn(func(_ interface{}, _, _ string) *redis.IntCmd {
				return redis.NewIntResult(1, nil)
			}).
			Times(1)

		if err := NewRedisQueueDriver(r).SetFailed(qName, taskId); err != nil {
			t.Errorf("Expected SetFailed() to set data, got error %v", err)
		}
	})
}

func TestRedisQueueDriver_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := NewMockRedisInterface(ctrl)

	qName := "test-queue"

	t.Run("it_should_return_error_when_it_fails_to_register_a_queue", func(t *testing.T) {
		expect := fmt.Errorf("failed to write")

		r.
			EXPECT().
			SAdd(gomock.Any(), "simple-queue:data:queue-list", "test-queue").
			DoAndReturn(func(_ interface{}, _, _ string) *redis.IntCmd {
				return redis.NewIntResult(0, fmt.Errorf("failed to write"))
			}).
			Times(1)

		if err := NewRedisQueueDriver(r).Register(qName); err == nil || err.Error() != expect.Error() {
			t.Errorf("Expected Register() to return error %v, got %v", expect, err)
		}
	})

	t.Run("it_should_register_queue", func(t *testing.T) {
		r.
			EXPECT().
			SAdd(gomock.Any(), "simple-queue:data:queue-list", "test-queue").
			DoAndReturn(func(_ interface{}, _, _ string) *redis.IntCmd {
				return redis.NewIntResult(1, nil)
			}).
			Times(1)

		if err := NewRedisQueueDriver(r).Register(qName); err != nil {
			t.Errorf("Expected Register() to register queue, got %v", err)
		}
	})
}

func TestRedisQueueDriver_GetStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := NewMockRedisInterface(ctrl)

	t.Run("it_should_return_usage_stats", func(t *testing.T) {
		r.
			EXPECT().
			SMembers(gomock.Any(), "simple-queue:data:queue-list").
			DoAndReturn(func(_ interface{}, _ string) *redis.StringSliceCmd {
				return redis.NewStringSliceResult([]string{"test-queue"}, nil)
			}).
			Times(1)

		r.
			EXPECT().
			Get(gomock.Any(), "simple-queue:data:test-queue:processed").
			DoAndReturn(func(_ interface{}, _ string) *redis.StringCmd {
				return redis.NewStringResult("103", nil)
			}).
			Times(1)

		fail1 := uuid.New().String()
		fail2 := uuid.New().String()

		r.
			EXPECT().
			SMembers(gomock.Any(), "simple-queue:data:test-queue:failed").
			DoAndReturn(func(_ interface{}, _ string) *redis.StringSliceCmd {
				return redis.NewStringSliceResult([]string{fail1, fail2}, nil)
			}).
			Times(1)

		expect := Stats{
			"test-queue": {
				Failed:    2,
				Processed: 103,
				FailedIDs: []string{fail1, fail2},
			},
		}

		if got, _ := NewRedisQueueDriver(r).GetStats(); !reflect.DeepEqual(got, &expect) {
			t.Errorf("Expected GetStats() to return %v, got %v", expect, got)
		}
	})
}
