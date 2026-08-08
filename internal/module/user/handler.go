package user

import (
	"github.com/google/uuid"

	"github.com/Eicap/EICAP-BANK/server/internal/enum"
	"github.com/Eicap/EICAP-BANK/server/internal/middleware"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	RegisterRoutes(v1 fiber.Router)
	Create(ctx fiber.Ctx) error
	Me(ctx fiber.Ctx) error
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
	user := v1.Group("/users")
	user.Use(middleware.JWT(h.jwtSecret))

	user.Post("/",
		middleware.RequirePermission(enum.Admin),
		h.Create,
	)

	user.Get("/me",
		h.Me,
	)

	user.Get("/:id",
		middleware.RequirePermission(enum.Admin),
		h.FindByID, // admin consulta cualquier usuario
	)

	user.Get("/",
		middleware.RequirePermission(enum.Admin),
		h.FindAll,
	)
}

func (h *handler) Create(c fiber.Ctx) error {
	var input Create
	if err := c.Bind().Body(&input); err != nil {
		return err
	}

	claims := middleware.User(c)

	if claims.Role == enum.Student {
		return response.Forbidden("El estudiante no puede crear un usuario")
	}

	if err := h.service.Create(c.Context(), &input); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response.OK("Usuario creado exitosamente", nil))
}

func (h *handler) Me(c fiber.Ctx) error {
	claims := middleware.User(c)
	result, err := h.service.FindByID(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return c.JSON(response.OK("Usuario encontrado exitosamente", result))
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
	return c.JSON(response.OK("Usuario encontrado exitosamente", result))
}

func (h *handler) FindAll(c fiber.Ctx) error {
	var filter UserFilter
	if err := c.Bind().Query(&filter); err != nil {
		return err
	}

	filter.SetDefaults()
	result, err := h.service.FindAll(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(response.OK("Usuarios encontrados exitosamente", result))
}
