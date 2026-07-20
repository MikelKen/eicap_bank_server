package main

import (
	bootstrap "github.com/Eicap/EICAP-BANK/server/internal/boottrap"
	"github.com/Eicap/EICAP-BANK/server/internal/config"
	"github.com/Eicap/EICAP-BANK/server/internal/database"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	app := fiber.New(fiber.Config{
		StructValidator: bootstrap.NewValidator(),
		ErrorHandler:    bootstrap.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowOrigins},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))

	container := bootstrap.NewContainer(db, cfg)
	bootstrap.SetupRoutes(app, container)

	app.Listen(":" + cfg.AppPort)
}
