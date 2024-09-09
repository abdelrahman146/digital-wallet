package main

import (
	"digital-wallet/internal/handler"
	"digital-wallet/internal/repository"
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/resources"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize the database connection
	db := resources.InitDB()

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			status, resp := api.NewErrorResponse(err)
			return ctx.Status(status).JSON(resp)
		},
	})

	// Middleware setup
	app.Use(recover.New())
	app.Use(healthcheck.New())
	app.Use(helmet.New())
	app.Use(idempotency.New())
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 10 * time.Second,
	}))
	app.Use(requestid.New())
	app.Get("/metrics", monitor.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "${time}: [${ip}:${port}] [${pid}] requestId:${locals:requestid} ${status} - ${method} ${path} ${latency}\n",
	}))

	// Define repositories
	repos := &repository.Repos{
		Account:      repository.NewAccountRepo(db),
		Transaction:  repository.NewTransactionRepo(db),
		Wallet:       repository.NewWalletRepo(db),
		User:         repository.NewUserRepo(db),
		Tier:         repository.NewTierRepo(db),
		ExchangeRate: repository.NewExchangeRateRepo(db),
	}

	// Define services
	services := &service.Services{
		Wallet:       service.NewWalletService(repos),
		Transaction:  service.NewTransactionService(repos),
		Account:      service.NewAccountService(repos),
		User:         service.NewUserService(repos),
		Tier:         service.NewTierService(repos),
		ExchangeRate: service.NewExchangeRateService(repos),
	}

	// Define versioned routes
	v1 := app.Group("/v1")

	handler.NewV1WalletHandler(v1, services)
	handler.NewV1AccountHandler(v1, services)
	handler.NewV1TransactionsHandler(v1, services)
	handler.NewV1ExchangeRateHandler(v1, services)
	handler.NewV1TierHandler(v1, services)
	handler.NewV1UserHandler(v1, services)
	handler.NewV1MainHandler(app, services)

	// Undefined route handler
	app.Use(func(c *fiber.Ctx) error {
		return errs.NewNotFoundError("Route not found", nil)
	})

	// Signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := app.Listen(":3401"); err != nil {
			logger.GetLogger().Panic("failed to start server", logger.Field("error", err))
		}
	}()

	// Wait for interrupt signal
	<-quit

	// Gracefully shutdown the server and close the database connection
	logger.GetLogger().Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.GetLogger().Error("Error shutting down server", logger.Field("error", err))
	}
	resources.CloseDB(db)
	logger.GetLogger().Info("Database connection closed")
}
