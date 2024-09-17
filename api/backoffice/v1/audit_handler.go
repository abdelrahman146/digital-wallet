package backofficev1

import (
	"github.com/abdelrahman146/digital-wallet/internal/service"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
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

// GetTableAuditLogs retrieves all audit logs of a table
// @Summary Get all audit logs of a table
// @Description Get all audit logs of a table
// @Tags Audit
// @Param table path string true "Table"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.Audit}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/audit/table/{table} [get]
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

// GetRecordAuditLogs retrieves all audit logs of a record
// @Summary Get all audit logs of a record
// @Description Get all audit logs of a record
// @Tags Audit
// @Param table path string true "Table"
// @Param recordId path string true "Record ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.Audit}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/audit/record/{table}/{recordId} [get]
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

// GetActorAuditLogs retrieves all audit logs of an actor
// @Summary Get all audit logs of an actor
// @Description Get all audit logs of an actor
// @Tags Audit
// @Param actor path string true "Actor"
// @Param actorId path string true "Actor ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.Audit}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/audit/actor/{actor}/{actorId} [get]
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
