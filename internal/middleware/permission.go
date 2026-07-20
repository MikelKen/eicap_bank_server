package middleware

import (
	"slices"

	"github.com/Eicap/EICAP-BANK/server/internal/enum"
	"github.com/gofiber/fiber/v3"
)

func RequirePermission(roles ...enum.Permission) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims := User(c)
		if claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		if slices.Contains(roles, claims.Role) {
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

}
