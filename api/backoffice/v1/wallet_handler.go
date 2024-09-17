package backofficev1

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"math"
)

type walletHandler struct {
	services *service.Services
}

func NewWalletHandler(appGroup fiber.Router, services *service.Services) {
	handler := &walletHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *walletHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("wallets")
	group.Post("/", h.CreateWallet)
	group.Get("/", h.GetWallets)
	group.Get("/:walletId/check-integrity", h.CheckWalletIntegrity)
	group.Get("/:walletId", h.GetWalletByID)
	group.Put("/:walletId", h.UpdateWallet)
	group.Delete("/:walletId", h.DeleteWallet)

}

// CreateWallet creates a new wallet
// @Summary Create a new wallet
// @Description Create a wallet based on the provided request
// @Tags Wallet
// @Accept json
// @Produce json
// @Param wallet body service.CreateWalletRequest true "Create Wallet Request"
// @Success 201 {object} api.SuccessResponse{result=model.Wallet}
// @Failure 400 {object} api.ErrorResponse
// @Router /wallets [post]
func (h *walletHandler) CreateWallet(c *fiber.Ctx) error {
	var req service.CreateWalletRequest
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	wallet, err := h.services.Wallet.CreateWallet(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(wallet))
}

func (h *walletHandler) GetWalletByID(c *fiber.Ctx) error {
	id := c.Params("walletId")
	wallet, err := h.services.Wallet.GetWalletByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(wallet))
}

func (h *walletHandler) GetWallets(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	wallets, err := h.services.Wallet.GetWallets(c.Context(), page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(wallets))
}

func (h *walletHandler) UpdateWallet(c *fiber.Ctx) error {
	id := c.Params("walletId")
	var req service.UpdateWalletRequest
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	wallet, err := h.services.Wallet.UpdateWallet(c.Context(), id, &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(wallet))
}

func (h *walletHandler) DeleteWallet(c *fiber.Ctx) error {
	id := c.Params("walletId")
	err := h.services.Wallet.DeleteWallet(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}

func (h *walletHandler) CheckWalletIntegrity(c *fiber.Ctx) error {
	id := c.Params("walletId")
	accountsSum, err := h.services.Wallet.GetAccountsSum(c.Context(), id)
	if err != nil {
		return err
	}
	transactionsSum, err := h.services.Wallet.GetTransactionsSum(c.Context(), id)
	if err != nil {
		return err
	}
	diff := math.Abs(float64(accountsSum - transactionsSum))
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(fiber.Map{
		"accountsSum":     accountsSum,
		"transactionsSum": transactionsSum,
		"diff":            diff,
	}))
}
