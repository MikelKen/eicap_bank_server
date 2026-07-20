package response

import "github.com/gofiber/fiber/v3"

func BadRequest(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusBadRequest, msg)
}

func Unauthorized(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusUnauthorized, msg)
}

func Forbidden(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusForbidden, msg)
}

func NotFound(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusNotFound, msg)
}

func Conflict(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusConflict, msg)
}

func InternalServerError(msg string) *fiber.Error {
	return fiber.NewError(fiber.StatusInternalServerError, msg)
}
