package backofficev1

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type transactionHandler struct {
	services *service.Services
}

func NewTransactionHandler(appGroup fiber.Router, services *service.Services) {
	handler := &transactionHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *transactionHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("wallets/:walletId/transactions")
	group.Post("/", h.CreateTransaction)
	group.Post("/exchange", h.CreateExchangeTransaction)
	group.Get("/", h.GetTransactions)
	group.Get("/sum", h.GetTransactionsSum)
}

func (h *transactionHandler) GetTransactions(c *fiber.Ctx) error {
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

func (h *transactionHandler) GetTransactionsSum(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	sum, err := h.services.Transaction.GetWalletTransactionSum(c.Context(), walletId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}

func (h *transactionHandler) CreateTransaction(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	var req struct {
		service.TransactionRequest
		AccountId string `json:"accountId,omitempty" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	transaction, err := h.services.Transaction.CreateTransaction(c.Context(), walletId, req.AccountId, &req.TransactionRequest)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(transaction))
}

func (h *transactionHandler) CreateExchangeTransaction(c *fiber.Ctx) error {
	var req struct {
		FromWalletID string `json:"fromWalletId,omitempty" validate:"required"`
		ToWalletID   string `json:"toWalletId,omitempty" validate:"required"`
		UserID       string `json:"userId,omitempty" validate:"required"`
		Amount       uint64 `json:"amount,omitempty" validate:"required,gt=0"`
	}
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	// Validate request
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		api.GetLogger(c.Context()).Error("Invalid request", logger.Field("fields", fields))
		return errs.NewValidationError("Invalid request", "", fields)
	}

	exchangeResponse, err := h.services.Transaction.Exchange(c.Context(), req.FromWalletID, req.ToWalletID, req.UserID, req.Amount)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(exchangeResponse))
}
