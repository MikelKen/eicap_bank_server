package bootstrap

import "github.com/gofiber/fiber/v3"

func SetupRoutes(app *fiber.App, container *Container) {
	v1 := app.Group("/api/v1")

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Welcome to EICAP-BANK API")
	})

	for _, handler := range container.Handlers {
		handler.RegisterRoutes(v1)
	}
}
