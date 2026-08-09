package bootstrap

import (
	"github.com/Eicap/EICAP-BANK/server/internal/config"
	"github.com/Eicap/EICAP-BANK/server/internal/module/account"
	"github.com/Eicap/EICAP-BANK/server/internal/module/auth"
	bankoperation "github.com/Eicap/EICAP-BANK/server/internal/module/bank_operation"
	cashcount "github.com/Eicap/EICAP-BANK/server/internal/module/cash_count"
	cashsession "github.com/Eicap/EICAP-BANK/server/internal/module/cash_session"
	"github.com/Eicap/EICAP-BANK/server/internal/module/client"
	"github.com/Eicap/EICAP-BANK/server/internal/module/denomination"
	typeaccount "github.com/Eicap/EICAP-BANK/server/internal/module/type_account"
	typeoperation "github.com/Eicap/EICAP-BANK/server/internal/module/type_operation"
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

	typeOperationRepo := typeoperation.NewRepo(db)
	typeOperationService := typeoperation.NewService(typeOperationRepo)
	typeOperationHandler := typeoperation.NewHandler(typeOperationService, cfg.JWTSecret)

	accountRepo := account.NewRepo(db)
	accountService := account.NewService(accountRepo, clientRepo, typeAccountRepo)
	accountHandler := account.NewHandler(accountService, cfg.JWTSecret)

	denominationRepo := denomination.NewRepo(db)
	denominationService := denomination.NewService(denominationRepo)
	denominationHandler := denomination.NewHandler(denominationService, cfg.JWTSecret)

	cashCountRepo := cashcount.NewRepo(db)
	cashCountService := cashcount.NewService(cashCountRepo)
	cashCountHandler := cashcount.NewHandler(cashCountService, cfg.JWTSecret)

	cashSessionRepo := cashsession.NewRepo(db)
	cashSessionService := cashsession.NewService(cashSessionRepo, denominationRepo) // reusa el repo ya creado
	cashSessionHandler := cashsession.NewHandler(cashSessionService, cfg.JWTSecret)

	bankOperationRepo := bankoperation.NewRepo(db)
	bankOperationService := bankoperation.NewService(bankOperationRepo, accountRepo, typeOperationRepo, cashSessionRepo)
	bankOperationHandler := bankoperation.NewHandler(bankOperationService, cfg.JWTSecret)

	return &Container{
		Handlers: []RouterRegister{
			userHandler,
			authHandler,
			clientHandler,
			typeAccountHandler,
			accountHandler,
			typeOperationHandler,
			denominationHandler,
			cashCountHandler,
			cashSessionHandler,
			bankOperationHandler,
		},
	}
}
