package bootstrap

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	// 1. Errores de Validación (go-playground/validator)
	if valErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
		msgs := make(map[string]string)
		for _, e := range valErrors {
			msgs[e.Field()] = e.Tag()
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation_failed",
			"details": msgs,
		})
	}
	// 2. Errores genéricos de Fiber (404, etc)
	if fiberErr, ok := errors.AsType[*fiber.Error](err); ok {
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"error": fiberErr.Message,
		})
	}

	// 3. Error por defecto (500)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "internal_server_error: " + err.Error(),
	})
}
