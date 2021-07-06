package simpleq

import (
	"fmt"
	"time"
)

const (
	LogLevelError = iota
	LogLevelWarn
	LogLevelInfo
)

// LogLevel package global
var LogLevel = LogLevelError

const (
	infoColor    = "\033[1;34m%s\033[0m"
	warningColor = "\033[1;33m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
)

var logLevelText = map[int]string{
	LogLevelError: "ERROR",
	LogLevelWarn:  "WARNING",
	LogLevelInfo:  "INFO",
}

// Logger is a log interface
type Logger interface {
	Error(err error)
	Info(i interface{})
	Warn(w interface{})
}

// DefaultLogger implementation
type DefaultLogger struct{}

// Error log
func (dl *DefaultLogger) Error(err error) {
	dl.write(LogLevelError, errorColor, err)
}

// Info log
func (dl *DefaultLogger) Info(i interface{}) {
	dl.write(LogLevelInfo, infoColor, i)
}

// Warn log
func (dl *DefaultLogger) Warn(w interface{}) {
	dl.write(LogLevelWarn, warningColor, w)
}

func (dl *DefaultLogger) write(level int, color string, i interface{}) {
	if LogLevel >= level {
		var layout = "02-Jan-2006 15:04:05"

		fmt.Printf(color, fmt.Sprintf("%v [%v] %v\n", time.Now().Format(layout), logLevelText[level], i))
	}
}
