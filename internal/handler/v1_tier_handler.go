package handler

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type v1TierHandler struct {
	services *service.Services
}

func NewV1TierHandler(appGroup fiber.Router, services *service.Services) {
	group := appGroup.Group("/tiers")
	handler := &v1TierHandler{services: services}
	handler.Setup(group)
}

func (h *v1TierHandler) Setup(group fiber.Router) {
	group.Post("/", h.CreateTier)
	group.Get("/", h.GetTiers)
	group.Get("/:tierId", h.GetTierByID)
	group.Delete("/:tierId", h.DeleteTier)
}

func (h *v1TierHandler) CreateTier(c *fiber.Ctx) error {
	var req service.CreateTierRequest
	if err := c.BodyParser(&req); err != nil {
		logger.GetLogger().Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	tier, err := h.services.Tier.CreateTier(&req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(tier))
}

func (h *v1TierHandler) GetTierByID(c *fiber.Ctx) error {
	id := c.Params("tierId")
	tier, err := h.services.Tier.GetTierByID(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(tier))
}

func (h *v1TierHandler) GetTiers(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	tiers, err := h.services.Tier.GetTiers(page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(tiers))
}

func (h *v1TierHandler) DeleteTier(c *fiber.Ctx) error {
	id := c.Params("tierId")
	err := h.services.Tier.DeleteTier(id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}
