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
	appGroup.Get("/ping")
}

func (h *walletHandler) Ping(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("pong")
}
