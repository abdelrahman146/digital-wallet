package handler

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"github.com/gofiber/fiber/v2"
)

type v1TransactionsHandler struct {
	services *service.Services
}

func NewV1TransactionsHandler(appGroup fiber.Router, services *service.Services) {
	handler := &v1TransactionsHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *v1TransactionsHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("wallets/:walletId/transactions")
	group.Get("/", h.GetTransactions)
	group.Get("/sum", h.GetTransactionsSum)
}

func (h *v1TransactionsHandler) GetTransactions(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	walletId := c.Params("walletId")
	transactions, err := h.services.Transaction.GetWalletTransactions(c.Context(), walletId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}

func (h *v1TransactionsHandler) GetTransactionsSum(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	sum, err := h.services.Transaction.GetWalletTransactionSum(c.Context(), walletId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}
