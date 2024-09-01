package handler

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"github.com/gofiber/fiber/v2"
	"net/http"
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
	group := appGroup.Group("/transactions")
	group.Post("/deposit", h.Deposit)
	group.Post("/withdraw", h.Withdraw)
	group.Post("/refund", h.Refund)
	group.Post("/purchase", h.Purchase)
	group.Post("/transfer", h.Transfer)
	group.Get("/sum", h.GetTransactionsSum)
}

func (h *v1TransactionsHandler) Deposit(c *fiber.Ctx) error {
	var req service.DepositRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("invalid request", err)
	}
	transaction, err := h.services.Transaction.Deposit(&req, model.TransactionInitiatedByUser)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(api.NewSuccessResponse(transaction))
}

func (h *v1TransactionsHandler) Withdraw(c *fiber.Ctx) error {
	var req service.WithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("invalid request", err)
	}
	transaction, err := h.services.Transaction.Withdraw(&req, model.TransactionInitiatedByUser)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(api.NewSuccessResponse(transaction))
}

func (h *v1TransactionsHandler) Refund(c *fiber.Ctx) error {
	var req service.RefundRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("invalid request", err)
	}
	transaction, err := h.services.Transaction.Refund(&req, model.TransactionInitiatedByUser)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(api.NewSuccessResponse(transaction))
}

func (h *v1TransactionsHandler) Purchase(c *fiber.Ctx) error {
	var req service.PurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("invalid request", err)
	}
	transaction, err := h.services.Transaction.Purchase(&req, model.TransactionInitiatedByUser)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(api.NewSuccessResponse(transaction))
}

func (h *v1TransactionsHandler) Transfer(c *fiber.Ctx) error {
	var req service.TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("invalid request", err)
	}
	transaction, err := h.services.Transaction.Transfer(&req, model.TransactionInitiatedByUser)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(api.NewSuccessResponse(transaction))
}

func (h *v1TransactionsHandler) GetTransactionsSum(c *fiber.Ctx) error {
	sum, err := h.services.Transaction.GetTransactionsSum()
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(api.NewSuccessResponse(sum))
}
