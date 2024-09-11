package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
	level  zapcore.Level
}

func newZapLogger(level zapcore.Level) (Logger, error) {
	logger := &zapLogger{level: level}
	if err := logger.Setup(); err != nil {
		return nil, err
	}
	return logger, nil
}

func (zl *zapLogger) Setup() error {
	config := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zl.level),
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
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, err := config.Build()
	if err != nil {
		return err
	}
	zl.logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

func (zl *zapLogger) Debug(msg string, fields ...F) {
	var zapFields []zap.Field
	for i := 0; i < len(fields); i++ {
		zapFields = append(zapFields, zap.Any(fields[i].Key, fields[i].Value))
	}
	zl.logger.Debug(msg, zapFields...)
}

func (zl *zapLogger) Info(msg string, fields ...F) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	zl.logger.Info(msg, zapFields...)
}

func (zl *zapLogger) Warn(msg string, fields ...F) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	zl.logger.Warn(msg, zapFields...)
}

func (zl *zapLogger) Error(msg string, fields ...F) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	zl.logger.Error(msg, zapFields...)
}

func (zl *zapLogger) Panic(msg string, fields ...F) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	zl.logger.Panic(msg, zapFields...)
}

func (zl *zapLogger) AddRequestID(key, value string) {
	_ = zl.Setup()
	zl.logger = zl.logger.With(zap.String(key, value))
}
