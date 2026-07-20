package pkg

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type CookieConfig struct {
	HttpOnly bool
	Secure   bool
	SameSite string
	Path     string
	Domain   string
}

func DefaultCookieConfig(isProduction bool) CookieConfig {
	return CookieConfig{
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
	}
}

func (cfg CookieConfig) NewAuthCookie(token string, expiration time.Duration) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: cfg.HttpOnly,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Path:     cfg.Path,
		Domain:   cfg.Domain,
		MaxAge:   int(expiration.Seconds()),
	}
}

func ClearAuthCookie(isProduction bool) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "token",
		Value:    "",
		HTTPOnly: true,
		Secure:   isProduction,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1, // Set MaxAge to -1 to delete the cookie
	}
}
