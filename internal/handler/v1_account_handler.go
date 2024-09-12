package handler

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type v1AccountHandler struct {
	services *service.Services
}

func NewV1AccountHandler(appGroup fiber.Router, services *service.Services) {
	handler := &v1AccountHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *v1AccountHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("wallets/:walletId/accounts")
	group.Get("/", h.GetAccounts)
	group.Post("/", h.CreateAccount)
	group.Get("/sum", h.GetAccountsSum)
	group.Get("/:accountId", h.GetAccountByID)
	group.Delete("/:accountId", h.DeleteAccount)
	group.Get("/:accountId/transactions", h.GetAccountTransactionsByID)
	group.Post("/:accountId/transactions", h.CreateTransaction)
	group.Post("/:accountId/transactions/sum", h.GetAccountTransactionsSum)
}

func (h *v1AccountHandler) GetAccounts(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	walletId := c.Params("walletId")
	if err != nil {
		return err
	}
	accounts, err := h.services.Account.GetAccounts(c.Context(), walletId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(accounts))
}

func (h *v1AccountHandler) GetAccountsSum(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	sum, err := h.services.Account.GetAccountsSum(c.Context(), walletId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}

func (h *v1AccountHandler) GetAccountByID(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	id := c.Params("accountId")
	account, err := h.services.Account.GetAccountByID(c.Context(), walletId, id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(account))
}

func (h *v1AccountHandler) CreateAccount(c *fiber.Ctx) error {
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

func (h *v1AccountHandler) GetAccountTransactionsByID(c *fiber.Ctx) error {
	id := c.Params("accountId")
	walletId := c.Params("walletId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	_, err = h.services.Account.GetAccountByID(c.Context(), walletId, id)
	if err != nil {
		return err
	}
	transactions, err := h.services.Transaction.GetTransactionsByAccountID(c.Context(), walletId, id, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}

func (h *v1AccountHandler) DeleteAccount(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	id := c.Params("accountId")
	err := h.services.Account.DeleteAccount(c.Context(), walletId, id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}

func (h *v1AccountHandler) GetAccountTransactionsByAccountID(c *fiber.Ctx) error {
	accountId := c.Params("accountId")
	walletId := c.Params("walletId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	transactions, err := h.services.Transaction.GetTransactionsByAccountID(c.Context(), walletId, accountId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}

func (h *v1AccountHandler) CreateTransaction(c *fiber.Ctx) error {
	accountId := c.Params("accountId")
	walletId := c.Params("walletId")
	var req service.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	transaction, err := h.services.Transaction.CreateTransaction(c.Context(), walletId, accountId, model.TransactionActorTypeUser, "123", &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(transaction))
}

func (h *v1AccountHandler) GetAccountTransactionsSum(c *fiber.Ctx) error {
	accountId := c.Params("accountId")
	walletId := c.Params("walletId")
	sum, err := h.services.Transaction.GetTransactionsSumByAccountID(c.Context(), walletId, accountId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}
