package response

import "github.com/gofiber/fiber/v3"

func OK(msg string, data any) fiber.Map {
	return fiber.Map{"message": msg, "data": data}
}
