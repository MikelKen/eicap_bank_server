package middleware

import (
	"log"

	"github.com/Eicap/EICAP-BANK/server/pkg"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

func JWT(jwtSecret string) fiber.Handler {
	log.Println("jwtSecret: ", jwtSecret)
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecret)},
		Extractor:  extractors.FromCookie("token"),
		Claims:     &pkg.Claims{},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	})

}

func User(c fiber.Ctx) *pkg.Claims {
	t := jwtware.FromContext(c)
	if t == nil {
		return nil
	}

	claims, ok := t.Claims.(*pkg.Claims)
	if !ok {
		return nil
	}
	return claims
}
