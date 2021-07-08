### SimpleQ

Simple queue implementation

---

#### Features

- [x] Queue
- [x] Requeue
- [x] Basic Logging 
- [ ] Detailed Logging
- [ ] Schedule
- [ ] Reschedule
- [ ] Delete
- [ ] Simple Stats UI

##### Usage Example

```go
package main

import (
	"github.com/go-redis/redis"
	"simpleq"
)

var redisClient = redis.NewClient(&redis.Options{})

// Task implementation
type Task struct{}

// Run is called every time task is executed
func (t *Task) Run(c simpleq.Context) error {
	// @TO-DO run things

	return nil
}

// Fail called when there is a failure
func (t *Task) Fail(err error) {
	// @TO-DO log your errors
}

func main() {
	driver := simpleq.NewRedisQueueDriver(redisClient)
	simpleq.Init(driver, &simpleq.DefaultLogger{})

	// 5 => number of concurrent handlers (should not set very high)
	q, err := simpleq.NewQueue("queue-name", 5)

	if err != nil {
		panic(err)
	}

	// push to the defined queue
	if err := q.Push(simpleq.NewMessage([]byte(uuid.New().String()))); err != nil {
		panic(err)
	}

	// specify which task to execute 
	q.OnExec(new(Task))

	sigC := make(chan os.Signal)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
	<-sigC

	// graceful shutdown...
	q.Stop()
	<-q.StopC
}
```
