package backofficev1

import (
	"github.com/abdelrahman146/digital-wallet/internal/service"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	services *service.Services
}

func NewUserHandler(appGroup fiber.Router, services *service.Services) {
	handler := &userHandler{services: services}
	handler.Setup(appGroup)
}

func (h *userHandler) Setup(appGroup fiber.Router) {
	group := appGroup.Group("/users")
	group.Post("/", h.CreateUser)
	group.Get("/", h.GetUsers)
	group.Get("/:userId", h.GetUserByID)
	group.Put("/:userId/tier", h.SetUserTier)
	group.Delete("/:userId", h.DeleteUser)
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a user based on the provided request
// @Tags User
// @Accept json
// @Produce json
// @Param user body service.CreateUserRequest true "Create User Request"
// @Success 201 {object} api.SuccessResponse{result=model.User}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/users [post]
func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	var req service.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	user, err := h.services.User.CreateUser(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(api.NewSuccessResponse(user))
}

// GetUserByID retrieves a user by its ID
// @Summary Get a user by its ID
// @Description Get a user by its ID
// @Tags User
// @Param userId path string true "User ID"
// @Success 200 {object} api.SuccessResponse{result=model.User}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/users/{userId} [get]
func (h *userHandler) GetUserByID(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user, err := h.services.User.GetUserByID(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(user))
}

// SetUserTier sets the tier of a user
// @Summary Set the tier of a user
// @Description Set the tier of a user
// @Tags User
// @Param userId path string true "User ID"
// @Param tierId path string true "Tier ID"
// @Success 202 {object} api.SuccessResponse{result=model.User}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/users/{userId}/tier/{tierId} [put]
func (h *userHandler) SetUserTier(c *fiber.Ctx) error {
	userId := c.Params("userId")
	tierId := c.Params("tierId")
	user, err := h.services.User.SetUserTier(c.Context(), userId, tierId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(user))
}

// GetUsers retrieves a list of users
// @Summary Get a list of users
// @Description Get a list of users
// @Tags User
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} api.SuccessResponse{result=[]model.User}
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/users [get]
func (h *userHandler) GetUsers(c *fiber.Ctx) error {
	page, limit, err := api.GetPageAndLimit(c)
	if err != nil {
		return err
	}
	users, err := h.services.User.GetUsers(c.Context(), page, limit)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(users))
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Description Delete a user
// @Tags User
// @Param userId path string true "User ID"
// @Success 202 {object} api.SuccessResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /backoffice/users/{userId} [delete]
func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	err := h.services.User.DeleteUser(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}
