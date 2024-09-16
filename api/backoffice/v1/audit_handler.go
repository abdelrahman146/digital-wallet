package backofficev1

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"github.com/gofiber/fiber/v2"
)

type auditHandler struct {
	services *service.Services
}

func NewAuditHandler(appGroup fiber.Router, services *service.Services) {
	handler := &auditHandler{services: services}
	handler.Setup(appGroup)
}

func (h *auditHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("audit")
	group.Get("/table/:table", h.GetTableAuditLogs)
	group.Get("/record/:table/:recordId", h.GetRecordAuditLogs)
	group.Get("/actor/:actor/:actorId", h.GetActorAuditLogs)
}

func (h *auditHandler) GetTableAuditLogs(c *fiber.Ctx) error {
	table := c.Params("table")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	audits, err := h.services.Audit.GetTableAuditLogs(c.Context(), table, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(audits))
}

func (h *auditHandler) GetRecordAuditLogs(c *fiber.Ctx) error {
	table := c.Params("table")
	recordId := c.Params("recordId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	audits, err := h.services.Audit.GetRecordAuditLogs(c.Context(), table, recordId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(audits))
}

func (h *auditHandler) GetActorAuditLogs(c *fiber.Ctx) error {
	actor := c.Params("actor")
	actorId := c.Params("actorId")
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	audits, err := h.services.Audit.GetActorAuditLogs(c.Context(), actor, actorId, page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(audits))
}
