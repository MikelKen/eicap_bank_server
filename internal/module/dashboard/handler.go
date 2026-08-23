package dashboard

import (
	"github.com/Eicap/EICAP-BANK/server/internal/middleware"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	RegisterRoutes(v1 fiber.Router)
	Summary(ctx fiber.Ctx) error
}

type handler struct {
	service   Service
	jwtSecret string
}

func NewHandler(service Service, jwtSecret string) Handler {
	return &handler{service: service, jwtSecret: jwtSecret}
}

func (h *handler) RegisterRoutes(v1 fiber.Router) {
	g := v1.Group("/dashboard")
	g.Use(middleware.JWT(h.jwtSecret))

	g.Get("/summary", h.Summary)
}

func (h *handler) Summary(c fiber.Ctx) error {
	claims := middleware.User(c)
	if claims == nil {
		return response.Unauthorized("Acceso no autorizado")
	}

	result, err := h.service.Summary(c.Context(), claims.UserID, claims.Role)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Resumen del panel encontrado exitosamente", result))
}
