package typeaccount

import (
	"github.com/Eicap/EICAP-BANK/server/internal/enum"
	"github.com/Eicap/EICAP-BANK/server/internal/middleware"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler interface {
	RegisterRoutes(v1 fiber.Router)
	Create(ctx fiber.Ctx) error
	Update(ctx fiber.Ctx) error
	FindByID(ctx fiber.Ctx) error
	FindAll(ctx fiber.Ctx) error
	Delete(ctx fiber.Ctx) error
}

type handler struct {
	service   Service
	jwtSecret string
}

func NewHandler(service Service, jwtSecret string) Handler {
	return &handler{service: service, jwtSecret: jwtSecret}
}

func (h *handler) RegisterRoutes(v1 fiber.Router) {
	typeAccount := v1.Group("/type-accounts")
	typeAccount.Use(middleware.JWT(h.jwtSecret))

	typeAccount.Post("/",
		middleware.RequirePermission(enum.Admin),
		h.Create,
	)

	typeAccount.Put("/:id",
		middleware.RequirePermission(enum.Admin),
		h.Update,
	)

	typeAccount.Get("/:id", h.FindByID)

	typeAccount.Get("/", h.FindAll)

	typeAccount.Delete("/:id",
		middleware.RequirePermission(enum.Admin),
		h.Delete,
	)
}

func (h *handler) Create(c fiber.Ctx) error {
	var input Create
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	if err := h.service.Create(c.Context(), &input); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response.OK("Tipo de cuenta creado exitosamente", nil))
}

func (h *handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest("ID inválido")
	}

	var input Update
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	if err := h.service.Update(c.Context(), id, &input); err != nil {
		return err
	}
	return c.JSON(response.OK("Tipo de cuenta actualizado exitosamente", nil))
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
	return c.JSON(response.OK("Tipo de cuenta encontrado exitosamente", result))
}

func (h *handler) FindAll(c fiber.Ctx) error {
	var filter TypeAccountFilter
	if err := c.Bind().Query(&filter); err != nil {
		return err
	}
	filter.SetDefaults()

	result, err := h.service.FindAll(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Tipos de cuenta encontrados exitosamente", result))
}

func (h *handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest("ID inválido")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(response.OK("Tipo de cuenta eliminado exitosamente", nil))
}
