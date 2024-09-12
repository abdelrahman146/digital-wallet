package handler

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type v1ExchangeRateHandler struct {
	services *service.Services
}

func NewV1ExchangeRateHandler(appGroup fiber.Router, services *service.Services) {
	handler := &v1ExchangeRateHandler{
		services: services,
	}
	handler.Setup(appGroup)
}

func (h *v1ExchangeRateHandler) Setup(group fiber.Router) {
	group.Post("wallets/exchange-rates", h.CreateExchangeRate)
	group.Get("/exchange-rates", h.GetExchangeRates)
	group.Get("/wallets/:walletId/exchange-rates", h.GetExchangeRatesByWalletID)
	group.Post("wallets/:walletId/exchange", h.Exchange)
	group.Put("wallets/exchange-rates/:exchangeRateId", h.UpdateExchangeRate)
	group.Delete("wallets/exchange-rates/:exchangeRateId", h.DeleteExchangeRate)
}

func (h *v1ExchangeRateHandler) CreateExchangeRate(c *fiber.Ctx) error {
	var req service.CreateExchangeRateRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	exchangeRate, err := h.services.ExchangeRate.CreateExchangeRate(&req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(exchangeRate))
}

func (h *v1ExchangeRateHandler) GetExchangeRates(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	exchangeRates, err := h.services.ExchangeRate.GetExchangeRates(page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(exchangeRates))
}

func (h *v1ExchangeRateHandler) GetExchangeRatesByWalletID(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	walletId := c.Params("walletId")
	exchangeRates, err := h.services.ExchangeRate.GetExchangeRatesByWalletID(walletId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(exchangeRates))
}

func (h *v1ExchangeRateHandler) UpdateExchangeRate(c *fiber.Ctx) error {
	exchangeRateId := c.Params("exchangeRateId")
	var req struct {
		ExchangeRate decimal.Decimal `json:"exchangeRate" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		logger.GetLogger().Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	exchangeRate, err := h.services.ExchangeRate.UpdateExchangeRate(exchangeRateId, req.ExchangeRate)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(exchangeRate))
}

func (h *v1ExchangeRateHandler) DeleteExchangeRate(c *fiber.Ctx) error {
	exchangeRateId := c.Params("exchangeRateId")
	err := h.services.ExchangeRate.DeleteExchangeRate(exchangeRateId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}

func (h *v1ExchangeRateHandler) Exchange(c *fiber.Ctx) error {
	fromWalletId := c.Params("walletId")
	var req struct {
		ToWalletID string `json:"toWalletId,omitempty" validate:"required"`
		UserID     string `json:"userId,omitempty" validate:"required"`
		Amount     uint64 `json:"amount,omitempty" validate:"required,gt=0"`
	}

	if err := c.BodyParser(&req); err != nil {
		logger.GetLogger().Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}

	// Validate request
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		logger.GetLogger().Error("Invalid request", logger.Field("fields", fields))
		return errs.NewValidationError("Invalid request", "", fields)
	}

	exchangeResponse, err := h.services.ExchangeRate.Exchange(fromWalletId, req.ToWalletID, req.UserID, model.TransactionActorTypeUser, req.UserID, req.Amount)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(exchangeResponse))
}
