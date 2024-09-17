package backofficev1

import (
	"github.com/abdelrahman146/digital-wallet/internal/service"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/abdelrahman146/digital-wallet/pkg/validator"
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

// GetWalletAccounts retrieves all accounts of a wallet
// @Summary Get all accounts of a wallet
// @Description Get all accounts of a wallet
// @Tags Account
// @Param walletId path string true "Wallet ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.Account}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/wallets/{walletId}/accounts [get]
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

// GetWalletAccountsSum retrieves the sum of all accounts of a wallet
// @Summary Get the sum of all accounts of a wallet
// @Description Get the sum of all accounts of a wallet
// @Tags Account
// @Param walletId path string true "Wallet ID"
// @Success 200 {object} api.SuccessResponse{result=float64}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/wallets/{walletId}/accounts/sum [get]
func (h *accountHandler) GetWalletAccountsSum(c *fiber.Ctx) error {
	walletId := c.Params("walletId")
	sum, err := h.services.Account.GetWalletAccountsSum(c.Context(), walletId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}

// GetAccountByID retrieves an account by its ID
// @Summary Get an account by its ID
// @Description Get an account by its ID
// @Tags Account
// @Param walletId path string true "Wallet ID"
// @Param accountId path string true "Account ID"
// @Success 200 {object} api.SuccessResponse{result=model.Account}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/wallets/{walletId}/accounts/{accountId} [get]
func (h *accountHandler) GetAccountByID(c *fiber.Ctx) error {
	id := c.Params("accountId")
	account, err := h.services.Account.GetAccount(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(account))
}

// CreateAccount creates an account
// @Summary Create an account
// @Description Create an account based on the provided request
// @Tags Account
// @Accept json
// @Produce json
// @Param walletId path string true "Wallet ID"
// @Param account body object true "Create Account Request"
// @Success 201 {object} api.SuccessResponse{result=model.Account}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/wallets/{walletId}/accounts [post]
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

// GetAccountTransactionsByID retrieves all transactions of an account
// @Summary Get all transactions of an account
// @Description Get all transactions of an account
// @Tags Account
// @Param walletId path string true "Wallet ID"
// @Param accountId path string true "Account ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.Transaction}
// @Failure 400 {object} api.ErrorResponse
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

// DeleteAccount deletes an account
// @Summary Delete an account
// @Description Delete an account
// @Tags Account
// @Param walletId path string true "Wallet ID"
// @Param accountId path string true "Account ID"
// @Success 202 {object} api.SuccessResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/wallets/{walletId}/accounts/{accountId} [delete]
func (h *accountHandler) DeleteAccount(c *fiber.Ctx) error {
	id := c.Params("accountId")
	err := h.services.Account.DeleteAccount(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}

// GetAccountTransactions retrieves all transactions of an account
// @Summary Get all transactions of an account
// @Description Get all transactions of an account
// @Tags Account
// @Param walletId path string true "Wallet ID"
// @Param accountId path string true "Account ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.Transaction}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/wallets/{walletId}/accounts/{accountId}/transactions [get]
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

// GetAccountTransactionsSum retrieves the sum of all transactions of an account
// @Summary Get the sum of all transactions of an account
// @Description Get the sum of all transactions of an account
// @Tags Account
// @Param walletId path string true "Wallet ID"
// @Param accountId path string true "Account ID"
// @Success 200 {object} api.SuccessResponse{result=float64}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/wallets/{walletId}/accounts/{accountId}/transactions/sum [get]
func (h *accountHandler) GetAccountTransactionsSum(c *fiber.Ctx) error {
	accountId := c.Params("accountId")
	walletId := c.Params("walletId")
	sum, err := h.services.Transaction.GetAccountTransactionSum(c.Context(), walletId, accountId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(sum))
}
