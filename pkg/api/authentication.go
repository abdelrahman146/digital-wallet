package api

import "github.com/gofiber/fiber/v2"

func AdminAuthenticationMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get jwt token
		// validate jwt token
		// get user id from jwt token
		ctx.Locals("actorId", "admin_id")
		return ctx.Next()
	}
}

func UserAuthenticationMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// get jwt token
		// validate jwt token
		// get user id from jwt token
		ctx.Locals("actorId", "user_id")
		return ctx.Next()
	}
}
