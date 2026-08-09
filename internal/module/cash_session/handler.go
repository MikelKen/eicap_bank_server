package cashsession

import (
	"github.com/Eicap/EICAP-BANK/server/internal/middleware"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler interface {
	RegisterRoutes(v1 fiber.Router)
	Open(ctx fiber.Ctx) error
	Close(ctx fiber.Ctx) error
	FindByID(ctx fiber.Ctx) error
	FindAll(ctx fiber.Ctx) error
	FindMyOpen(ctx fiber.Ctx) error
}

type handler struct {
	service   Service
	jwtSecret string
}

func NewHandler(service Service, jwtSecret string) Handler {
	return &handler{service: service, jwtSecret: jwtSecret}
}

func (h *handler) RegisterRoutes(v1 fiber.Router) {
	cashSession := v1.Group("/cash-sessions")
	cashSession.Use(middleware.JWT(h.jwtSecret))

	cashSession.Post("/open", h.Open)
	cashSession.Put("/:id/close", h.Close)
	cashSession.Get("/mine/open", h.FindMyOpen)
	cashSession.Get("/:id", h.FindByID)
	cashSession.Get("/", h.FindAll)
}

func (h *handler) Open(c fiber.Ctx) error {
	var input Open
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	claims := middleware.User(c)
	if claims == nil {
		return response.Unauthorized("Acceso no autorizado")
	}

	if err := h.service.Open(c.Context(), claims.UserID, &input); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response.OK("Caja abierta exitosamente", nil))
}

func (h *handler) Close(c fiber.Ctx) error {
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest("ID de sesión inválido")
	}

	var input Close
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	claims := middleware.User(c)
	if claims == nil {
		return response.Unauthorized("Acceso no autorizado")
	}

	if err := h.service.Close(c.Context(), claims.UserID, sessionID, &input); err != nil {
		return err
	}
	return c.JSON(response.OK("Caja cerrada exitosamente", nil))
}

func (h *handler) FindByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest("ID inválido")
	}

	result, err := h.service.FindByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Sesión de caja encontrada exitosamente", result))
}

func (h *handler) FindAll(c fiber.Ctx) error {
	var filter CashSessionFilter
	if err := c.Bind().Query(&filter); err != nil {
		return err
	}
	filter.SetDefaults()

	result, err := h.service.FindAll(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Sesiones de caja encontradas exitosamente", result))
}

func (h *handler) FindMyOpen(c fiber.Ctx) error {
	claims := middleware.User(c)
	if claims == nil {
		return response.Unauthorized("Acceso no autorizado")
	}

	result, err := h.service.FindOpenByUserID(c.Context(), claims.UserID)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Sesión de caja abierta encontrada", result))
}
