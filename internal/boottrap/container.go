package bootstrap

import (
	"github.com/Eicap/EICAP-BANK/server/internal/config"
	"github.com/Eicap/EICAP-BANK/server/internal/module/auth"
	"github.com/Eicap/EICAP-BANK/server/internal/module/user"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type RouterRegister interface {
	RegisterRoutes(fiber.Router)
}

type Container struct {
	Handlers []RouterRegister
}

func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, cfg.JWTSecret)

	authService := auth.NewService(userRepo, cfg)
	authHandler := auth.NewHandler(authService, cfg.AppEnv == "production")

	return &Container{
		Handlers: []RouterRegister{
			userHandler,
			authHandler,
		},
	}
}
