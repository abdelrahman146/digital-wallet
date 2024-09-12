package logger

import (
	"go.uber.org/zap/zapcore"
	"sync"
)

type F struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Debug(msg string, fields ...F)
	Info(msg string, fields ...F)
	Warn(msg string, fields ...F)
	Error(msg string, fields ...F)
	Panic(msg string, fields ...F)
}

var (
	logger Logger
	once   sync.Once
)

func initLogger(l Logger) {
	once.Do(func() {
		logger = l
	})
}

func GetLogger() Logger {
	if logger == nil {
		zap, err := NewZapLogger(zapcore.DebugLevel)
		if err != nil {
			panic(err)
		}
		initLogger(zap)
	}
	return logger
}

func Field(key string, value interface{}) F {
	return F{Key: key, Value: value}
}
