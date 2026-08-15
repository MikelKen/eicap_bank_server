package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	AppEnv       string
	DatabaseURL  string
	RunMigration bool
	RunSeeder    bool
	AllowOrigins string

	JWTSecret     string
	JWTExpiration string

	AdminPasswordOne string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(`.env file not found`)
	}
	log.Println("URL: ", os.Getenv("ALLOW_ORIGINS"))

	return &Config{
		AppPort:      os.Getenv("APP_PORT"),
		AppEnv:       os.Getenv("APP_ENV"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		RunMigration: os.Getenv("RUN_MIGRATION") == "true",
		RunSeeder:    os.Getenv("RUN_SEEDER") == "true",
		AllowOrigins: os.Getenv("ALLOW_ORIGINS"),

		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),

		AdminPasswordOne: os.Getenv("ADMIN_PASSWORD_ONE"),
	}
}
