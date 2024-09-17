package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
	level  zapcore.Level
}

func NewZapLogger(level zapcore.Level, fields ...F) (Logger, error) {
	zl := &zapLogger{level: level}
	logger, err := zl.Setup(fields...)
	zl.logger = logger
	if err != nil {
		return nil, err
	}
	return zl, nil
}

func (zl *zapLogger) Setup(fields ...F) (*zap.Logger, error) {
	config := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zl.level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:    "time",
			LevelKey:   "level",
			NameKey:    "logger",
			CallerKey:  "caller",
			MessageKey: "msg",
			//StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
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
	for _, field := range fields {
		logger = logger.With(zap.Any(field.Key, field.Value))
	}
	return logger, nil
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
