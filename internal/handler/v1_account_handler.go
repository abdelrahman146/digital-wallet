package handler

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type accountHandler struct {
	services *service.Services
}

func NewV1AccountHandler(appGroup fiber.Router, services *service.Services) {
	handler := &accountHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *accountHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("/accounts")
	group.Get("/", h.GetAccounts)
	group.Post("/", h.CreateAccount)
	group.Get("/sum", h.GetAccountsSum)
	group.Get("/user/:userId", h.GetAccountByUserID)
	group.Get("/user/:userId/transactions", h.GetAccountTransactionsByUserID)
	group.Get("/:accountId", h.GetAccountByID)
	group.Get("/:accountId/transactions", h.GetAccountTransactionsByID)
}

func (h *accountHandler) GetAccounts(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	accounts, err := h.services.Account.GetAccounts(page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(accounts))
}

func (h *accountHandler) GetAccountsSum(c *fiber.Ctx) error {
	sum, err := h.services.Account.GetAccountsSum()
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}

func (h *accountHandler) GetAccountByID(c *fiber.Ctx) error {
	id := c.Params("accountId")
	account, err := h.services.Account.GetAccountByID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(account))
}

func (h *accountHandler) GetAccountByUserID(c *fiber.Ctx) error {
	id := c.Params("userId")
	account, err := h.services.Account.GetAccountByUserID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(account))
}

func (h *accountHandler) CreateAccount(c *fiber.Ctx) error {
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

	account, err := h.services.Account.CreateAccount(req.UserID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(account))
}

func (h *accountHandler) GetAccountTransactionsByID(c *fiber.Ctx) error {
	id := c.Params("accountId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	_, err = h.services.Account.GetAccountByID(id)
	if err != nil {
		return err
	}
	transactions, err := h.services.Transaction.GetTransactionsByAccountID(id, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}

func (h *accountHandler) GetAccountTransactionsByUserID(c *fiber.Ctx) error {
	id := c.Params("userId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	account, err := h.services.Account.GetAccountByUserID(id)
	if err != nil {
		return err
	}
	transactions, err := h.services.Transaction.GetTransactionsByAccountID(account.ID, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(transactions))
}
