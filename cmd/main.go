package main

import (
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"net/http"
)

func main() {
	//db := resources.InitDB()
	app := fiber.New()
	app.Use(recover.New())
	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(api.NewSuccessResponse("pong"))
	})
	err := app.Listen(":3401")
	if err != nil {
		logger.GetLogger().Panic("failed to start server", logger.Field("error", err))
	}
	// TODO CLOSE CONNECTIONS
}
