package logger

import "sync"

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

func InitLogger(l Logger) {
	once.Do(func() {
		logger = l
	})
}

func GetLogger() Logger {
	return logger
}

func Field(key string, value interface{}) F {
	return F{Key: key, Value: value}
}
