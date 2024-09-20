package backofficev1

import (
	"github.com/abdelrahman146/digital-wallet/internal/service"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type triggerHandler struct {
	services *service.Services
}

func NewTriggerHandler(appGroup fiber.Router, services *service.Services) {
	handler := &triggerHandler{services: services}
	handler.Setup(appGroup)
}

func (h *triggerHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("triggers")
	group.Post("/", h.CreateTrigger)
	group.Put("/:triggerId", h.UpdateTrigger)
	group.Delete("/:triggerId", h.DeleteTrigger)
	group.Get("/:triggerId", h.GetTrigger)
	group.Get("/", h.GetTriggers)
}

// CreateTrigger creates a new trigger
// @Summary Create a new trigger
// @Description Create a trigger based on the provided request
// @Tags Trigger
// @Accept json
// @Produce json
// @Param trigger body service.CreateTriggerRequest true "Create Trigger Request"
// @Success 201 {object} api.SuccessResponse{result=model.Trigger}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/triggers [post]
func (h *triggerHandler) CreateTrigger(c *fiber.Ctx) error {
	var req service.CreateTriggerRequest
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request")
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	trigger, err := h.services.Trigger.CreateTrigger(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(trigger))
}

// UpdateTrigger updates a trigger
// @Summary Update a trigger
// @Description Update a trigger based on the provided request
// @Tags Trigger
// @Accept json
// @Produce json
// @Param triggerId path string true "Trigger ID"
// @Param trigger body service.UpdateTriggerRequest true "Update Trigger Request"
// @Success 200 {object} api.SuccessResponse{result=model.Trigger}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/triggers/{triggerId} [put]
func (h *triggerHandler) UpdateTrigger(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("triggerId"), 10, 64)
	if err != nil {
		api.GetLogger(c.Context()).Error("Invalid trigger ID", logger.Field("triggerId", c.Params("triggerId")))
		return errs.NewBadRequestError("Invalid trigger ID", "INVALID_TRIGGER_ID", err)
	}
	var req service.UpdateTriggerRequest
	if err := c.BodyParser(&req); err != nil {
		api.GetLogger(c.Context()).Error("Invalid body request")
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	trigger, err := h.services.Trigger.UpdateTrigger(c.Context(), id, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(trigger))
}

// DeleteTrigger deletes a trigger
// @Summary Delete a trigger
// @Description Delete a trigger based on the provided request
// @Tags Trigger
// @Param triggerId path string true "Trigger ID"
// @Success 202 {object} api.SuccessResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/triggers/{triggerId} [delete]
func (h *triggerHandler) DeleteTrigger(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("triggerId"), 10, 64)
	if err != nil {
		api.GetLogger(c.Context()).Error("Invalid trigger ID", logger.Field("triggerId", c.Params("triggerId")))
		return errs.NewBadRequestError("Invalid trigger ID", "INVALID_TRIGGER_ID", err)
	}
	err = h.services.Trigger.DeleteTrigger(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}

// GetTrigger retrieves a trigger by its ID
// @Summary Get a trigger by its ID
// @Description Get a trigger by its ID
// @Tags Trigger
// @Param triggerId path string true "Trigger ID"
// @Success 200 {object} api.SuccessResponse{result=model.Trigger}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/triggers/{triggerId} [get]
func (h *triggerHandler) GetTrigger(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("triggerId"), 10, 64)
	if err != nil {
		api.GetLogger(c.Context()).Error("Invalid trigger ID", logger.Field("triggerId", c.Params("triggerId")))
		return errs.NewBadRequestError("Invalid trigger ID", "INVALID_TRIGGER_ID", err)
	}
	trigger, err := h.services.Trigger.GetTrigger(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(trigger))
}

// GetTriggers retrieves a list of triggers
// @Summary Get a list of triggers
// @Description Get a list of triggers
// @Tags Trigger
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.Trigger}
// @Failure 400 {object} api.ErrorResponse
func (h *triggerHandler) GetTriggers(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	triggers, err := h.services.Trigger.ListTriggers(c.Context(), page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(triggers))
}
