package handler

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"github.com/gofiber/fiber/v2"
)

type v1walletHandler struct {
	services *service.Services
}

func NewV1WalletHandler(appGroup fiber.Router, services *service.Services) {
	handler := &v1walletHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *v1walletHandler) Setup(group fiber.Router) {
	group.Post("/", h.CreateWallet)
	group.Get("/", h.GetWallets)
	group.Get("/:walletId", h.GetWalletByID)
	group.Put("/:walletId", h.UpdateWallet)
	group.Delete("/:walletId", h.DeleteWallet)
}

func (h *v1walletHandler) CreateWallet(c *fiber.Ctx) error {
	var req service.CreateWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("invalid request", err)
	}
	wallet, err := h.services.Wallet.CreateWallet(&req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(wallet))
}

func (h *v1walletHandler) GetWalletByID(c *fiber.Ctx) error {
	id := c.Params("walletId")
	wallet, err := h.services.Wallet.GetWalletByID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(wallet))
}

func (h *v1walletHandler) GetWallets(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	wallets, err := h.services.Wallet.GetWallets(page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(wallets))
}

func (h *v1walletHandler) UpdateWallet(c *fiber.Ctx) error {
	id := c.Params("walletId")
	var req service.UpdateWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("invalid request", err)
	}
	wallet, err := h.services.Wallet.UpdateWallet(id, &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(wallet))
}

func (h *v1walletHandler) DeleteWallet(c *fiber.Ctx) error {
	id := c.Params("walletId")
	err := h.services.Wallet.DeleteWallet(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusNoContent).JSON(api.NewSuccessResponse(nil))
}
