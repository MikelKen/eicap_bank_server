package bankoperation

import (
	"github.com/Eicap/EICAP-BANK/server/internal/middleware"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler interface {
	RegisterRoutes(v1 fiber.Router)
	Create(ctx fiber.Ctx) error
	FindByID(ctx fiber.Ctx) error
	FindAll(ctx fiber.Ctx) error
}

type handler struct {
	service   Service
	jwtSecret string
}

func NewHandler(service Service, jwtSecret string) Handler {
	return &handler{service: service, jwtSecret: jwtSecret}
}

func (h *handler) RegisterRoutes(v1 fiber.Router) {
	g := v1.Group("/bank-operations")
	g.Use(middleware.JWT(h.jwtSecret))

	g.Post("/", h.Create)
	g.Get("/:id", h.FindByID)
	g.Get("/", h.FindAll)
}

func (h *handler) Create(c fiber.Ctx) error {
	var input Create
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	claims := middleware.User(c)
	if claims == nil {
		return response.Unauthorized("Acceso no autorizado")
	}

	if err := h.service.Create(c.Context(), claims.UserID, &input); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response.OK("Operación registrada exitosamente", nil))
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
	return c.JSON(response.OK("Operación encontrada exitosamente", result))
}

func (h *handler) FindAll(c fiber.Ctx) error {
	var filter BankOperationFilter
	if err := c.Bind().Query(&filter); err != nil {
		return err
	}
	filter.SetDefaults()

	result, err := h.service.FindAll(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Operaciones encontradas exitosamente", result))
}
