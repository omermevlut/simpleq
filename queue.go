package simpleq

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

const (
	queuePrefix = "simple-queue:data"
)

var (
	driver Driver
	logger Logger
)

// Init initializes simple queue with a given driver implementation
func Init(d Driver, l Logger) {
	driver = d
	logger = l
}

// Queueable is a queue interface
type Queueable interface {
	Push(t Context) error
	OnExec(task Task)
	Requeue(t Context) error
	Stop()
}

// NewQueue returns a pointer to a new Queue instance
func NewQueue(name string, workers int8) (*Queue, error) {
	if err := driver.Register(name); err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Initialized queue %s", name))

	return &Queue{
		Workers:  workers,
		Name:     name,
		StopC:    make(chan struct{}),
		stopExec: make(chan struct{}),
	}, nil
}

// Queue is an instance queue handler
// use DelayedQueue for delayed (scheduled) messages
// use NewQueue() factory function instead of manually initializing
type Queue struct {
	Workers  int8
	Name     string
	StopC    chan struct{}
	stopExec chan struct{}
}

// Push to queue
func (q *Queue) Push(c Context) error {
	c.SetID()

	d, err := c.Marshal()

	if err != nil {
		return err
	}

	return driver.Write(q.getActiveName(), d)
}

// OnExec is triggered when there is a new message in th queue
func (q *Queue) OnExec(task Task) {
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			select {
			case <-q.stopExec:
				close(q.StopC)
			case <-ticker.C:
				q.read(task)
			}
		}
	}()
}

// Requeue pushes the task back in into queue until max attempts reached
func (q *Queue) Requeue(c Context) error {
	c.NewAttempt()

	if c.GetAttempts() > c.GetMaxAttempts() {
		return fmt.Errorf("max attempts reached for %s:%s", q.Name, c.GetID())
	}

	return q.Push(c)
}

// Stop queue from being executed
func (q *Queue) Stop() {
	q.stopExec <- struct{}{}
}

func (q *Queue) read(task Task) {
	if d, err := driver.Read(q.getActiveName()); err != nil && err != redis.Nil {
		logger.Warn(err)
	} else if len(d) > 0 {
		var m Message
		var qName = fmt.Sprintf("%s:%s", queuePrefix, q.Name)

		if err := json.Unmarshal(d, &m); err != nil {
			logger.Warn(err)

			return
		}

		logger.Info(fmt.Sprintf("[Processing] queue %v, task ID: %v", q.Name, m.GetID()))
		if err := task.Run(&m); err != nil {
			task.Fail(err)
			_ = driver.SetFailed(qName, m.GetID())
			logger.Warn(fmt.Sprintf("[Failed] queue %v, task ID: %v", q.Name, m.GetID()))
		} else {
			_ = driver.SetProcessed(qName)
			logger.Info(fmt.Sprintf("[Processed] queue %v, task ID: %v", q.Name, m.GetID()))
		}
	}
}

func (q *Queue) getActiveName() string {
	return fmt.Sprintf("%s:active:%s", queuePrefix, q.Name)
}
