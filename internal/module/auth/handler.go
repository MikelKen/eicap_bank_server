package auth

import (
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg"
	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	RegisterRoutes(v1 fiber.Router)
	Login(c fiber.Ctx) error
}

type handler struct {
	service      Service
	isProduction bool
}

func NewHandler(service Service, isProduction bool) Handler {
	return &handler{service: service, isProduction: isProduction}
}

func (h *handler) RegisterRoutes(v1 fiber.Router) {
	auth := v1.Group("/auth")
	auth.Post("/login", h.Login)
}

func (h *handler) Login(c fiber.Ctx) error {
	var input Login
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	token, expiration, userResponse, err := h.service.Login(&input)
	if err != nil {
		return err
	}

	c.Cookie(pkg.DefaultCookieConfig(h.isProduction).NewAuthCookie(token, expiration))

	return c.Status(fiber.StatusOK).JSON(response.OK("Login exitoso", userResponse))
}
