package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger  *zap.Logger
	service string
}

func NewZapLogger(level zapcore.Level, service string) (Logger, error) {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	logger = logger.With(zap.String("service", service))
	return &ZapLogger{logger: logger, service: service}, nil
}

func (zl *ZapLogger) Debug(msg string, fields ...F) {
	var zapFields []zap.Field
	for i := 0; i < len(fields); i++ {
		zapFields = append(zapFields, zap.Any(fields[i].Key, fields[i].Value))
	}
	zl.logger.Debug(msg, zapFields...)
}

func (zl *ZapLogger) Info(msg string, fields ...F) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	zl.logger.Info(msg, zapFields...)
}

func (zl *ZapLogger) Warn(msg string, fields ...F) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	zl.logger.Warn(msg, zapFields...)
}

func (zl *ZapLogger) Error(msg string, fields ...F) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	zl.logger.Error(msg, zapFields...)
}

func (zl *ZapLogger) Panic(msg string, fields ...F) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	zl.logger.Panic(msg, zapFields...)
}
