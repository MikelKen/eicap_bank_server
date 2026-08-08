package database

import (
	"log"
	"path/filepath"
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/config"
	"github.com/Eicap/EICAP-BANK/server/internal/database/seed"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	maxRetries := 10

	for i := range maxRetries {
		if cfg.DatabaseURL == "" {
			log.Fatal("DATABASE_URL not set in .env file")
		}

		var err error
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
		if err == nil {
			log.Printf("Database connected successfully after %d attempt(s)", i+1)

			if cfg.RunMigration {
				migrate(db)
			}

			if cfg.RunSeeder {
				log.Printf("Running seeder...")
				seed.Run(db, cfg.AdminPasswordOne)
			}
			return db
		}
		log.Printf("Failed to connet to database, retryig (%d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Could not connect to database after retries")
	return nil

}

func migrate(db *gorm.DB) {
	dir, err := filepath.Abs("./internal/database/migration")
	if err != nil {
		log.Fatalf("Failed to get migration dir: %v", err)
	}
	log.Println("Resolved migration dir:", dir) // 👈 agrega esto
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get raw db: %v", err)
	}

	if err := goose.Up(sqlDB, dir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")
}
