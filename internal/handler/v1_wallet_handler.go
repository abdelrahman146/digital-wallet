package handler

import (
	"digital-wallet/internal/service"
	"github.com/gofiber/fiber/v2"
)

type walletHandler struct {
	services *service.Services
}

func NewV1WalletHandler(appGroup fiber.Router, services *service.Services) {
	handler := &walletHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *walletHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("/wallet")
	group.Get("/:userId", h.GetWalletByID)
}

func (h *walletHandler) GetWalletByID(c *fiber.Ctx) error {
	id := c.Params("userId")
	wallet, err := h.services.Wallet.GetWalletByUserID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(wallet)
}
