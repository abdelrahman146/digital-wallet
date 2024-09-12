package api

import (
	"context"
	"digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap/zapcore"
)

func CreateAppContext() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestId := ctx.Locals("requestid").(string)
		l, err := logger.NewZapLogger(zapcore.DebugLevel, logger.Field("requestId", requestId))
		if err != nil {
			return err
		}
		ctx.Locals("logger", l)
		return ctx.Next()
	}
}

func GetLogger(ctx context.Context) logger.Logger {
	return ctx.Value("logger").(logger.Logger)
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value("requestId").(string)
}
