package client

import (
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
	FindAllByUserID(ctx fiber.Ctx) error
	GetClient(ctx fiber.Ctx) error
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
	client := v1.Group("/clients")
	client.Use(middleware.JWT(h.jwtSecret))

	client.Post("/", h.Create)

	client.Get("/mine", h.FindAllByUserID)

	client.Put("/:id", h.Update)

	client.Get("/:id", h.FindByID)

	client.Get("/", h.FindAll)

	client.Get("/search/:data", h.GetClient)

	client.Delete("/:id", h.Delete)
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
	return c.Status(fiber.StatusCreated).JSON(response.OK("Cliente creado exitosamente", nil))
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
	return c.JSON(response.OK("Cliente actualizado exitosamente", nil))
}

func (h *handler) FindAllByUserID(c fiber.Ctx) error {
	claims := middleware.User(c)
	if claims == nil {
		return response.Unauthorized("Acceso no autorizado")
	}

	var filter ClientFilter
	if err := c.Bind().Query(&filter); err != nil {
		return err
	}
	filter.SetDefaults()

	result, err := h.service.FindAllByUserID(c.Context(), claims.UserID, filter)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Clientes encontrados exitosamente", result))
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
	return c.JSON(response.OK("Cliente encontrado exitosamente", result))
}

func (h *handler) FindAll(c fiber.Ctx) error {
	result, total, err := h.service.FindAll(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Clientes encontrados exitosamente", fiber.Map{
		"items": result,
		"total": total,
	}))
}

func (h *handler) GetClient(c fiber.Ctx) error {
	data := c.Params("data")
	if data == "" {
		return response.BadRequest("Dato de búsqueda inválido")
	}

	result, err := h.service.GetClient(c.Context(), data)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Cliente encontrado exitosamente", result))
}

func (h *handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest("ID inválido")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(response.OK("Cliente eliminado exitosamente", nil))
}
