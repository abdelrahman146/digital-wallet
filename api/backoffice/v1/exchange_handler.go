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

func (h *exchangeRateHandler) DeleteExchangeRate(c *fiber.Ctx) error {
	exchangeRateId := c.Params("exchangeRateId")
	err := h.services.ExchangeRate.DeleteExchangeRate(c.Context(), exchangeRateId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}
