package simpleq

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

type TaskImpl struct{}

func (ti *TaskImpl) Run(c Context) error {
	return nil
}

func (ti *TaskImpl) Fail(err error) {
	fmt.Println(err)
}

func TestNewQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := NewMockDriver(ctrl)

	qName := "test-queue"

	Init(d, &DefaultLogger{})

	t.Run("it_should_return_error_when_it_fails_to_initialize_a_queue", func(t *testing.T) {
		expect := fmt.Errorf("failed to register")

		d.
			EXPECT().
			Register(qName).
			DoAndReturn(func(_ string) error {
				return fmt.Errorf("failed to register")
			}).
			Times(1)

		if _, err := NewQueue(qName, 1); err == nil || err.Error() != expect.Error() {
			t.Errorf("Expected NewQueue() to return errror %v, got %v", expect, err)
		}
	})

	t.Run("it_should_return_new_queue_instance", func(t *testing.T) {
		d.
			EXPECT().
			Register(qName).
			DoAndReturn(func(_ string) error {
				return nil
			}).
			Times(1)

		expect := Queue{
			Workers:  2,
			Name:     "test-queue",
			StopC:    nil,
			stopExec: nil,
		}

		if got, err := NewQueue(qName, 2); err != nil {
			t.Errorf("Expected NewQueue() to reutrn a new instance, got error %v", err)
		} else {
			// Delete channels
			got.StopC = nil
			got.stopExec = nil

			if !reflect.DeepEqual(&expect, got) {
				t.Errorf("Expected NewQueue() to retutn %v, got %v", expect, got)
			}

		}
	})
}

func TestQueue_Push(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := NewMockContext(ctrl)
	d := NewMockDriver(ctrl)

	Init(d, &DefaultLogger{})

	queue := Queue{
		Workers:  1,
		Name:     "test-queue",
		StopC:    make(chan struct{}),
		stopExec: make(chan struct{}),
	}

	t.Run("it_should_return_error_when_it_fails_to_marshal_context", func(t *testing.T) {
		expect := fmt.Errorf("failed to marshal")

		c.EXPECT().SetID().Times(1)

		c.
			EXPECT().
			Marshal().
			DoAndReturn(func() ([]byte, error) {
				return nil, fmt.Errorf("failed to marshal")
			}).
			Times(1)

		if err := queue.Push(c); err == nil || err.Error() != expect.Error() {
			t.Errorf("Expected Push() to return error %v, got %v", expect, err)
		}
	})

	t.Run("it_should_return_error_when_it_fails_to_save_data", func(t *testing.T) {
		expect := fmt.Errorf("failed to write")

		c.EXPECT().SetID().Times(1)

		c.
			EXPECT().
			Marshal().
			DoAndReturn(func() ([]byte, error) {
				return []byte("test-data"), nil
			}).
			Times(1)

		d.
			EXPECT().
			Write("simple-queue:data:active:test-queue", []byte("test-data")).
			DoAndReturn(func(_ string, _ []byte) error {
				return fmt.Errorf("failed to write")
			}).
			Times(1)

		if err := queue.Push(c); err == nil || err.Error() != expect.Error() {
			t.Errorf("Expected Push() to return error %v, got %v", expect, err)
		}
	})
}

func TestQueue_read(t *testing.T) {
	ctrl := gomock.NewController(t)
	d := NewMockDriver(ctrl)
	dl := NewMockLogger(ctrl)
	task := NewMockTask(ctrl)

	queue := Queue{
		Workers:  3,
		Name:     "test-queue",
		StopC:    make(chan struct{}),
		stopExec: make(chan struct{}),
	}

	Init(d, dl)

	t.Run("it_should_log_error_when_it_fails_to_read_data", func(t *testing.T) {
		d.
			EXPECT().
			Read("simple-queue:data:active:test-queue").
			DoAndReturn(func(_ string) ([]byte, error) {
				return nil, fmt.Errorf("failed to read")
			}).
			Times(1)

		dl.EXPECT().Warn(fmt.Errorf("failed to read")).Times(1)

		queue.read(new(TaskImpl))
	})

	t.Run("it_should_log_warning_when_it_fails_to_unmarshal_data", func(t *testing.T) {
		d.
			EXPECT().
			Read("simple-queue:data:active:test-queue").
			DoAndReturn(func(_ string) ([]byte, error) {
				return []byte("{"), nil
			}).
			Times(1)

		dl.EXPECT().Warn(gomock.Any()).Times(1)

		queue.read(new(TaskImpl))
	})

	t.Run("it_should_call_fail_when_task_run_returns_an_error", func(t *testing.T) {
		d.
			EXPECT().
			Read("simple-queue:data:active:test-queue").
			DoAndReturn(func(_ string) ([]byte, error) {
				return []byte("{}"), nil
			}).
			Times(1)

		task.
			EXPECT().
			Run(gomock.Any()).
			DoAndReturn(func(_ interface{}) error {
				return fmt.Errorf("failed to run")
			}).
			Times(1)

		task.EXPECT().Fail(fmt.Errorf("failed to run")).Times(1)
		dl.EXPECT().Info(gomock.Any()).Times(2)

		d.EXPECT().SetFailed("simple-queue:data:test-queue", gomock.Any()).Times(1)
		queue.read(task)
	})

	t.Run("it_should_call_run", func(t *testing.T) {
		d.
			EXPECT().
			Read("simple-queue:data:active:test-queue").
			DoAndReturn(func(_ string) ([]byte, error) {
				return []byte("{}"), nil
			}).
			Times(1)

		task.
			EXPECT().
			Run(gomock.Any()).
			DoAndReturn(func(_ interface{}) error {
				return nil
			}).
			Times(1)

		dl.EXPECT().Info(gomock.Any()).Times(2)

		d.EXPECT().SetProcessed("simple-queue:data:test-queue").Times(1)
		queue.read(task)
	})
}

func TestQueue_Requeue(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := NewMockContext(ctrl)
	d := NewMockDriver(ctrl)
	lg := NewMockLogger(ctrl)

	Init(d, lg)

	queue := Queue{
		Workers:  1,
		Name:     "test-queue",
		StopC:    make(chan struct{}),
		stopExec: make(chan struct{}),
	}

	t.Run("it_should_return_error_when_max_attempts_reached", func(t *testing.T) {
		c.EXPECT().NewAttempt().Times(1)
		c.EXPECT().GetAttempts().DoAndReturn(func() int { return 5 }).Times(1)
		c.EXPECT().GetMaxAttempts().DoAndReturn(func() int { return 4 }).Times(1)
		c.EXPECT().GetID().DoAndReturn(func() string { return "test-id" }).Times(1)

		expect := fmt.Errorf("max attempts reached for %s:test-id", "test-queue")

		if err := queue.Requeue(c); err == nil || err.Error() != expect.Error() {
			t.Errorf("Expected Requeue() to return error %v, got %v", expect, err)
		}
	})

	t.Run("it_should_requeue", func(t *testing.T) {
		c.EXPECT().NewAttempt().Times(1)
		c.EXPECT().GetAttempts().DoAndReturn(func() int { return 1 }).Times(1)
		c.EXPECT().GetMaxAttempts().DoAndReturn(func() int { return 5 }).Times(1)

		c.EXPECT().SetID().Times(1)

		c.
			EXPECT().
			Marshal().
			DoAndReturn(func() ([]byte, error) {
				return []byte("test-data"), nil
			}).
			Times(1)

		d.
			EXPECT().
			Write("simple-queue:data:active:test-queue", []byte("test-data")).
			DoAndReturn(func(_ string, _ []byte) error {
				return nil
			}).
			Times(1)

		if err := queue.Requeue(c); err != nil {
			t.Errorf("Expected Requeue() to requeue, got error %v", err)
		}
	})
}

func TestQueue_Stop(t *testing.T) {
	t.Run("it_should_send_a_stop_signal", func(t *testing.T) {
		queue := Queue{
			Workers:  1,
			Name:     "test-queue",
			StopC:    make(chan struct{}),
			stopExec: make(chan struct{}),
		}

		queue.OnExec(new(TaskImpl))
		queue.Stop()

		expect := struct{}{}
		got := <-queue.StopC

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("Expected StopC to receive a stop signal")
		}
	})
}
