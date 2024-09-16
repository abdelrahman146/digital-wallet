package backofficev1

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App, services *service.Services) {
	group := app.Group("api/v1/backoffice/")
	group.Use(api.AdminAuthenticationMiddleware())
	group.Use(api.CreateAppContext(api.AppActorAdmin))
	NewAuditHandler(group, services)
	NewAccountHandler(group, services)
	NewExchangeRateHandler(group, services)
	NewTierHandler(group, services)
	NewUserHandler(group, services)
	NewWalletHandler(group, services)
	NewTransactionHandler(group, services)
}
