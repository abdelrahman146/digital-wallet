package handler

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type v1MainHandler struct {
	services *service.Services
}

func NewV1MainHandler(appGroup fiber.Router, services *service.Services) {
	handler := &v1MainHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *v1MainHandler) Setup(appGroup fiber.Router) {
	appGroup.Get("/check-integrity", h.CheckIntegrity)
}

func (h *v1MainHandler) CheckIntegrity(c *fiber.Ctx) error {
	transactionsSum, err := h.services.Transaction.GetTransactionsSum()
	if err != nil {
		return err
	}
	walletsSum, err := h.services.Account.GetAccountsSum()
	if err != nil {
		return err
	}
	diff := transactionsSum.Sub(walletsSum).Abs()
	return c.Status(http.StatusOK).JSON(api.NewSuccessResponse(fiber.Map{
		"transactionsSum": transactionsSum,
		"walletsSum":      walletsSum,
		"diff":            diff,
	}))
}
