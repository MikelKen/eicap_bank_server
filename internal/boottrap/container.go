package bootstrap

import (
	"github.com/Eicap/EICAP-BANK/server/internal/config"
	"github.com/Eicap/EICAP-BANK/server/internal/module/account"
	"github.com/Eicap/EICAP-BANK/server/internal/module/auth"
	"github.com/Eicap/EICAP-BANK/server/internal/module/client"
	typeaccount "github.com/Eicap/EICAP-BANK/server/internal/module/type_account"
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

	clientRepo := client.NewRepo(db)
	clientService := client.NewService(clientRepo)
	clientHandler := client.NewHandler(clientService, cfg.JWTSecret)

	typeAccountRepo := typeaccount.NewRepo(db)
	typeAccountService := typeaccount.NewService(typeAccountRepo)
	typeAccountHandler := typeaccount.NewHandler(typeAccountService, cfg.JWTSecret)

	accountRepo := account.NewRepo(db)
	accountService := account.NewService(accountRepo, clientRepo, typeAccountRepo) // reutiliza los repos que ya creaste
	accountHandler := account.NewHandler(accountService, cfg.JWTSecret)

	return &Container{
		Handlers: []RouterRegister{
			userHandler,
			authHandler,
			clientHandler,
			typeAccountHandler,
			accountHandler,
		},
	}
}
