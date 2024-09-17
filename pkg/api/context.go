package api

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap/zapcore"
)

func CreateAppContext(actor string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestId := ctx.Locals("requestid").(string)
		actorId := ctx.Locals("actorId").(string)
		l, err := logger.NewZapLogger(zapcore.DebugLevel, logger.Field("requestId", requestId), logger.Field("actor", actor), logger.Field("actorId", actorId))
		if err != nil {
			return err
		}
		ctx.Locals("logger", l)
		ctx.Locals("actor", actor)
		return ctx.Next()
	}
}

func GetLogger(ctx context.Context) logger.Logger {
	return ctx.Value("logger").(logger.Logger)
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value("requestId").(string)
}

func GetActorID(ctx context.Context) string {
	return ctx.Value("actorId").(string)
}

func GetActor(ctx context.Context) string {
	return ctx.Value("actor").(string)
}
