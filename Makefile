SHELL := /bin/bash

mock-gen-all:
	mockgen -destination=driver_mock.go -package="simpleq" -source=driver.go
	mockgen -destination=message_mock.go -package="simpleq" -source=message.go
	mockgen -destination=queue_mock.go -package="simpleq" -source=queue.go
	mockgen -destination=redis_mock.go -package="simpleq" -source=redis.go
	mockgen -destination=log_mock.go -package="simpleq" -source=log.go
