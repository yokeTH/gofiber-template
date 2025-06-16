package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func CtxRequestIDInjector() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Locals("requestid")
		if idStr, ok := id.(string); ok {
			// nolint:staticcheck
			ctx := context.WithValue(c.UserContext(), "requestid", idStr)
			c.SetUserContext(ctx)
		}

		return c.Next()
	}
}
