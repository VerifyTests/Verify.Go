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

func (l *stdLogger) Error(err error) {
	if err == nil {
		return
	}
	_ = l.info.Output(2, err.Error())
}

func (l *stdLogger) Info(format string, args ...interface{}) {
	_ = l.info.Output(2, fmt.Sprintf(format, args...))
}
