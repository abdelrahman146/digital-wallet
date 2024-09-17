package main

import (
	backofficev1 "github.com/abdelrahman146/digital-wallet/api/backoffice/v1"
	_ "github.com/abdelrahman146/digital-wallet/docs"
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	"github.com/abdelrahman146/digital-wallet/internal/service"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/abdelrahman146/digital-wallet/pkg/resources"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/swaggo/fiber-swagger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Digital Wallet API
// @version 1.0
// @description This is the Digital Wallet API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name Abdel Rahman Hussein
// @contact.url https://github.com/abdelrahman146

// @host localhost:3401
// @BasePath /api/v1
func main() {
	// Initialize the database connection
	db := resources.InitDB()

	// CreateTransaction a new Fiber app
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
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "${time}: [${ip}:${port}] [${pid}] requestId:${locals:requestid} ${status} - ${method} ${path} ${latency}\n",
	}))

	// Define repositories
	repos := &repository.Repos{
		Audit:        repository.NewAuditRepo(db),
		Account:      repository.NewAccountRepo(db),
		Transaction:  repository.NewTransactionRepo(db),
		Wallet:       repository.NewWalletRepo(db),
		User:         repository.NewUserRepo(db),
		Tier:         repository.NewTierRepo(db),
		ExchangeRate: repository.NewExchangeRateRepo(db),
	}

	// Define services
	services := &service.Services{
		Audit:        service.NewAuditService(repos),
		Wallet:       service.NewWalletService(repos),
		Transaction:  service.NewTransactionService(repos),
		Account:      service.NewAccountService(repos),
		User:         service.NewUserService(repos),
		Tier:         service.NewTierService(repos),
		ExchangeRate: service.NewExchangeRateService(repos),
	}

	// Define routes
	backofficev1.New(app, services)

	// Undefined route handler
	app.Use(func(c *fiber.Ctx) error {
		logger.GetLogger().Info("Route not found", logger.Field("path", c.Path()))
		return errs.NewNotFoundError("Route not found", "", nil)
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
