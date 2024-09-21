package backofficev1

import (
	"github.com/abdelrahman146/digital-wallet/internal/service"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type programHandler struct {
	services *service.Services
}

func NewProgramHandler(appGroup fiber.Router, services *service.Services) {
	handler := &programHandler{services: services}
	handler.Setup(appGroup)
}

func (h *programHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("programs")
	group.Post("/", h.CreateProgram)
	group.Patch("/:programId", h.UpdateProgram)
	group.Delete("/:programId", h.DeleteProgram)
	group.Get("/:programId", h.GetProgram)
	group.Get("/", h.GetPrograms)
}

// CreateProgram creates a new program
// @Summary Create a new program
// @Description Create a program based on the provided request
// @Tags Program
// @Accept json
// @Produce json
// @Param program body service.CreateProgramRequest true "Create Program Request"
// @Success 201 {object} api.SuccessResponse{result=model.Program}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/programs [post]
func (h *programHandler) CreateProgram(c *fiber.Ctx) error {
	var req service.CreateProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	program, err := h.services.Program.CreateProgram(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(program))
}

// UpdateProgram updates a program
// @Summary Update a program
// @Description Update a program based on the provided request
// @Tags Program
// @Accept json
// @Produce json
// @Param programId path string true "Program ID"
// @Param program body service.UpdateProgramRequest true "Update Program Request"
// @Success 200 {object} api.SuccessResponse{result=model.Program}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/programs/{programId} [patch]
func (h *programHandler) UpdateProgram(c *fiber.Ctx) error {
	programID, err := strconv.ParseUint(c.Params("programId"), 10, 64)
	if err != nil {
		return errs.NewBadRequestError("Invalid program ID", "INVALID_PROGRAM_ID", err)
	}
	var req service.UpdateProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.NewBadRequestError("Invalid body request", "INVALID_BODY_REQUEST", err)
	}
	program, err := h.services.Program.UpdateProgram(c.Context(), programID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(program))
}

// DeleteProgram deletes a program
// @Summary Delete a program
// @Description Delete a program based on the provided request
// @Tags Program
// @Param programId path string true "Program ID"
// @Success 204
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/programs/{programId} [delete]
func (h *programHandler) DeleteProgram(c *fiber.Ctx) error {
	programID, err := strconv.ParseUint(c.Params("programId"), 10, 64)
	if err != nil {
		return errs.NewBadRequestError("Invalid program ID", "INVALID_PROGRAM_ID", err)
	}
	err = h.services.Program.DeleteProgram(c.Context(), programID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}

// GetProgram retrieves a program
// @Summary Get a program
// @Description Get a program based on the provided request
// @Tags Program
// @Param programId path string true "Program ID"
// @Success 200 {object} api.SuccessResponse{result=model.Program}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/programs/{programId} [get]
func (h *programHandler) GetProgram(c *fiber.Ctx) error {
	programID, err := strconv.ParseUint(c.Params("programId"), 10, 64)
	if err != nil {
		return errs.NewBadRequestError("Invalid program ID", "INVALID_PROGRAM_ID", err)
	}
	program, err := h.services.Program.GetProgram(c.Context(), programID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(program))
}

// GetPrograms retrieves a list of programs
// @Summary Get a list of programs
// @Description Get a list of programs based on the provided request
// @Tags Program
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=api.List[model.Program]}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/programs [get]
func (h *programHandler) GetPrograms(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	programs, err := h.services.Program.ListPrograms(c.Context(), page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(programs))
}
