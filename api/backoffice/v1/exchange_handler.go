package backofficev1

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type exchangeRateHandler struct {
	services *service.Services
}

func NewExchangeRateHandler(appGroup fiber.Router, services *service.Services) {
	handler := &exchangeRateHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *exchangeRateHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("exchange-rates")
	group.Post("/", h.CreateExchangeRate)
	group.Get("/", h.GetExchangeRates)
	group.Get("/wallets/:walletId", h.GetExchangeRatesByWalletID)
	group.Put("/wallets/:walletId", h.UpdateExchangeRate)
	group.Delete("/wallets/:walletId", h.DeleteExchangeRate)
}

// CreateExchangeRate creates a new exchange rate
// @Summary Create a new exchange rate
// @Description Create an exchange rate based on the provided request
// @Tags Exchange Rate
// @Accept json
// @Produce json
// @Param exchangeRate body service.CreateExchangeRateRequest true "Create Exchange Rate Request"
// @Success 201 {object} api.SuccessResponse{result=model.ExchangeRate}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/exchange-rates [post]
func (h *exchangeRateHandler) CreateExchangeRate(c *fiber.Ctx) error {
	var req service.CreateExchangeRateRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	exchangeRate, err := h.services.ExchangeRate.CreateExchangeRate(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(exchangeRate))
}

// GetExchangeRates retrieves all exchange rates
// @Summary Get all exchange rates
// @Description Get all exchange rates
// @Tags Exchange Rate
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.ExchangeRate}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/exchange-rates [get]
func (h *exchangeRateHandler) GetExchangeRates(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	exchangeRates, err := h.services.ExchangeRate.GetExchangeRates(c.Context(), page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(exchangeRates))
}

// GetExchangeRatesByWalletID retrieves all exchange rates of a wallet
// @Summary Get all exchange rates of a wallet
// @Description Get all exchange rates of a wallet
// @Tags Exchange Rate
// @Param walletId path string true "Wallet ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.ExchangeRate}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/exchange-rates/wallets/{walletId} [get]
func (h *exchangeRateHandler) GetExchangeRatesByWalletID(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	walletId := c.Params("walletId")
	exchangeRates, err := h.services.ExchangeRate.GetExchangeRatesByWalletID(c.Context(), walletId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(exchangeRates))
}

// UpdateExchangeRate updates an exchange rate
// @Summary Update an exchange rate
// @Description Update an exchange rate based on the provided request
// @Tags Exchange Rate
// @Accept json
// @Produce json
// @Param exchangeRateId path string true "Exchange Rate ID"
// @Param exchangeRate body object true "Update Exchange Rate Request"
// @Success 200 {object} api.SuccessResponse{result=model.ExchangeRate}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/exchange-rates/{exchangeRateId} [put]
func (h *exchangeRateHandler) UpdateExchangeRate(c *fiber.Ctx) error {
	exchangeRateId := c.Params("exchangeRateId")
	var req struct {
		ExchangeRate decimal.Decimal `json:"exchangeRate" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	exchangeRate, err := h.services.ExchangeRate.UpdateExchangeRate(c.Context(), exchangeRateId, req.ExchangeRate)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(exchangeRate))
}

// DeleteExchangeRate deletes an exchange rate
// @Summary Delete an exchange rate
// @Description Delete an exchange rate by its ID
// @Tags Exchange Rate
// @Produce json
// @Param exchangeRateId path string true "Exchange Rate ID"
// @Success 202 {object} api.SuccessResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/exchange-rates/{exchangeRateId} [delete]
func (h *exchangeRateHandler) DeleteExchangeRate(c *fiber.Ctx) error {
	exchangeRateId := c.Params("exchangeRateId")
	err := h.services.ExchangeRate.DeleteExchangeRate(c.Context(), exchangeRateId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}
