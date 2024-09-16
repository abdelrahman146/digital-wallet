package backofficev1

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type accountHandler struct {
	services *service.Services
}

func NewAccountHandler(appGroup fiber.Router, services *service.Services) {
	handler := &accountHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *accountHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("wallets/:walletId/accounts")
	group.Get("/", h.GetWalletAccounts)
	group.Post("/", h.CreateAccount)
	group.Get("/sum", h.GetWalletAccountsSum)
	group.Get("/:accountId", h.GetAccountByID)
	group.Delete("/:accountId", h.DeleteAccount)
	group.Get("/:accountId/transactions", h.GetAccountTransactionsByID)
	group.Post("/:accountId/transactions/sum", h.GetAccountTransactionsSum)
}

func (h *accountHandler) GetWalletAccounts(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	walletId := c.Params("walletId")
	accounts, err := h.services.Account.GetWalletAccounts(c.Context(), walletId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(accounts))
}

func (h *accountHandler) GetWalletAccountsSum(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	sum, err := h.services.Account.GetWalletAccountsSum(c.Context(), walletId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}

func (h *accountHandler) GetAccountByID(c *fiber.Ctx) error {
	id := c.Params("accountId")
	account, err := h.services.Account.GetAccount(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(account))
}

func (h *accountHandler) CreateAccount(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	var req struct {
		UserID string `json:"userId,omitempty" validate:"required"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}

	// Validate request
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		api.GetLogger(c.Context()).Error("Invalid request", logger.Field("fields", fields))
		return errs.NewValidationError("Invalid request", "INVALID_BODY_REQUEST", fields)
	}

	account, err := h.services.Account.CreateAccount(c.Context(), walletId, req.UserID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(account))
}

func (h *accountHandler) GetAccountTransactionsByID(c *fiber.Ctx) error {
	id := c.Params("accountId")
	walletId := c.Params("walletId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	transactions, err := h.services.Transaction.GetAccountTransactions(c.Context(), walletId, id, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}

func (h *accountHandler) DeleteAccount(c *fiber.Ctx) error {
	id := c.Params("accountId")
	err := h.services.Account.DeleteAccount(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}

func (h *accountHandler) GetAccountTransactions(c *fiber.Ctx) error {
	accountId := c.Params("accountId")
	walletId := c.Params("walletId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	transactions, err := h.services.Transaction.GetAccountTransactions(c.Context(), walletId, accountId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}

func (h *accountHandler) GetAccountTransactionsSum(c *fiber.Ctx) error {
	accountId := c.Params("accountId")
	walletId := c.Params("walletId")
	sum, err := h.services.Transaction.GetAccountTransactionSum(c.Context(), walletId, accountId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}
