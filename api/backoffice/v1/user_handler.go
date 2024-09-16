package backofficev1

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
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

func (h *userHandler) GetUserByID(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user, err := h.services.User.GetUserByID(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(user))
}

func (h *userHandler) SetUserTier(c *fiber.Ctx) error {
	userId := c.Params("userId")
	tierId := c.Params("tierId")
	user, err := h.services.User.SetUserTier(c.Context(), userId, tierId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(user))
}

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

func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	err := h.services.User.DeleteUser(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}
