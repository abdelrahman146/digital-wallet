package api

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap/zapcore"
)

func CreateAppContextMiddleware(actor string) fiber.Handler {
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

func CreateAppContext(ctx context.Context, actor, actorId, requestId string) context.Context {
	ctx = context.WithValue(ctx, "requestId", requestId)
	ctx = context.WithValue(ctx, "actorId", actorId)
	ctx = context.WithValue(ctx, "actor", actor)
	l, _ := logger.NewZapLogger(zapcore.DebugLevel, logger.Field("requestId", requestId), logger.Field("actor", actor), logger.Field("actorId", actorId))
	ctx = context.WithValue(ctx, "logger", l)
	return ctx
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
