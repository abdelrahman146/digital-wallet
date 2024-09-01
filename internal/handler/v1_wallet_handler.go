package handler

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/validator"
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
	group := appGroup.Group("/wallets")
	group.Get("/", h.GetWallets)
	group.Post("/", h.CreateWallet)
	group.Get("/sum", h.GetWalletsSum)
	group.Get("/user/:userId", h.GetWalletByUserID)
	group.Get("/user/:userId/transactions", h.GetWalletTransactionsByUserID)
	group.Get("/:walletId", h.GetWalletByID)
	group.Get("/:walletId/transactions", h.GetWalletTransactionsByID)
}

func (h *walletHandler) GetWallets(c *fiber.Ctx) error {
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

func (h *walletHandler) GetWalletsSum(c *fiber.Ctx) error {
	sum, err := h.services.Wallet.GetWalletsSum()
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}

func (h *walletHandler) GetWalletByID(c *fiber.Ctx) error {
	id := c.Params("walletId")
	wallet, err := h.services.Wallet.GetWalletByID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(wallet))
}

func (h *walletHandler) GetWalletByUserID(c *fiber.Ctx) error {
	id := c.Params("userId")
	wallet, err := h.services.Wallet.GetWalletByUserID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(wallet))
}

func (h *walletHandler) CreateWallet(c *fiber.Ctx) error {
	var req struct {
		UserID string `json:"userId,omitempty" validate:"required,uuid"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("Invalid Body Request", err)
	}

	// Validate request
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return errs.NewValidationError("Invalid Body Request", fields)
	}

	wallet, err := h.services.Wallet.CreateWallet(req.UserID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(wallet))
}

func (h *walletHandler) GetWalletTransactionsByID(c *fiber.Ctx) error {
	id := c.Params("walletId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	_, err = h.services.Wallet.GetWalletByID(id)
	if err != nil {
		return err
	}
	transactions, err := h.services.Transaction.GetTransactionsByWalletID(id, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}

func (h *walletHandler) GetWalletTransactionsByUserID(c *fiber.Ctx) error {
	id := c.Params("userId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	wallet, err := h.services.Wallet.GetWalletByUserID(id)
	if err != nil {
		return err
	}
	transactions, err := h.services.Transaction.GetTransactionsByWalletID(wallet.ID, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}
