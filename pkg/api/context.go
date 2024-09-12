package api

import (
	"context"
	"digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap/zapcore"
)

func CreateAppContext(ctx *fiber.Ctx) {
	requestId := ctx.Locals("requestid").(string)
	l, _ := logger.NewZapLogger(zapcore.DebugLevel, logger.Field("requestId", requestId))
	ctx.Locals("logger", l)
}

func GetLogger(ctx context.Context) logger.Logger {
	return ctx.Value("logger").(logger.Logger)
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value("requestId").(string)
}
