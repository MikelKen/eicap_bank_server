package cashcount

import (
	"github.com/Eicap/EICAP-BANK/server/internal/middleware"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler interface {
	RegisterRoutes(v1 fiber.Router)
	FindAllBySessionID(ctx fiber.Ctx) error
}

type handler struct {
	service   Service
	jwtSecret string
}

func NewHandler(service Service, jwtSecret string) Handler {
	return &handler{service: service, jwtSecret: jwtSecret}
}

func (h *handler) RegisterRoutes(v1 fiber.Router) {
	cashCount := v1.Group("/cash-counts")
	cashCount.Use(middleware.JWT(h.jwtSecret))

	cashCount.Get("/session/:sessionId", h.FindAllBySessionID)
}

func (h *handler) FindAllBySessionID(c fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Params("sessionId"))
	if err != nil {
		return response.BadRequest("ID de sesión inválido")
	}

	result, err := h.service.FindAllBySessionID(c.Context(), sessionID)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Conteo de caja encontrado exitosamente", result))
}
