package handler

import (
	"digital-wallet/internal/service"
	"digital-wallet/pkg/api"
	"github.com/gofiber/fiber/v2"
)

type v1UserHandler struct {
	services *service.Services
}

func NewV1UserHandler(appGroup fiber.Router, services *service.Services) {
	group := appGroup.Group("/users")
	handler := &v1UserHandler{services: services}
	handler.Setup(group)
}

func (h *v1UserHandler) Setup(r fiber.Router) {
	r.Post("/", h.CreateUser)
	r.Get("/", h.GetUsers)
	r.Get("/:userId", h.GetUserByID)
	r.Put("/:userId/tier", h.SetUserTier)
	r.Delete("/:userId", h.DeleteUser)
}

func (h *v1UserHandler) CreateUser(c *fiber.Ctx) error {
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

func (h *v1UserHandler) GetUserByID(c *fiber.Ctx) error {
	userId := c.Params("userId")
	user, err := h.services.User.GetUserByID(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(api.NewSuccessResponse(user))
}

func (h *v1UserHandler) SetUserTier(c *fiber.Ctx) error {
	userId := c.Params("userId")
	tierId := c.Params("tierId")
	user, err := h.services.User.SetUserTier(c.Context(), userId, tierId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(user))
}

func (h *v1UserHandler) GetUsers(c *fiber.Ctx) error {
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

func (h *v1UserHandler) DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	err := h.services.User.DeleteUser(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusAccepted).JSON(api.NewSuccessResponse(nil))
}
