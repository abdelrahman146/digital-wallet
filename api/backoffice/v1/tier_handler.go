package backofficev1

import (
	"github.com/abdelrahman146/digital-wallet/internal/service"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type tierHandler struct {
	services *service.Services
}

func NewTierHandler(appGroup fiber.Router, services *service.Services) {
	handler := &tierHandler{services: services}
	handler.Setup(appGroup)
}

func (h *tierHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("tiers")
	group.Post("/", h.CreateTier)
	group.Get("/", h.GetTiers)
	group.Get("/:tierId", h.GetTierByID)
	group.Delete("/:tierId", h.DeleteTier)
}

func (h *tierHandler) CreateTier(c *fiber.Ctx) error {
	var req service.CreateTierRequest
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request", logger.Field("error", err))
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	tier, err := h.services.Tier.CreateTier(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(tier))
}

func (h *tierHandler) GetTierByID(c *fiber.Ctx) error {
	id := c.Params("tierId")
	tier, err := h.services.Tier.GetTierByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(tier))
}

func (h *tierHandler) GetTiers(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	tiers, err := h.services.Tier.GetTiers(c.Context(), page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(tiers))
}

func (h *tierHandler) DeleteTier(c *fiber.Ctx) error {
	id := c.Params("tierId")
	err := h.services.Tier.DeleteTier(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}
