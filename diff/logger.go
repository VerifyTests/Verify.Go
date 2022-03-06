package diff

import (
	"fmt"
	"log"
	"os"
)

type stdLogger struct {
	info *log.Logger
	err  *log.Logger
}

// Logger interface to log info and error logs
type Logger interface {
	Error(err error)
	Info(format string, args ...interface{})
}

func newLogger(ns string) Logger {
	info := log.New(os.Stdout, "["+ns+"][INFO] ", log.LstdFlags)
	err := log.New(os.Stdout, "["+ns+"][ERROR] ", log.LstdFlags)
	return &stdLogger{
		info: info,
		err:  err,
	}
}

// Error logs an error message
func (l *stdLogger) Error(err error) {
	if err == nil {
		return
	}
	_ = l.info.Output(2, err.Error())
}

// Info logs an info message
func (l *stdLogger) Info(format string, args ...interface{}) {
	_ = l.info.Output(2, fmt.Sprintf(format, args...))
}
