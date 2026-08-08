package seed

import (
	"errors"
	"log"

	"github.com/Eicap/EICAP-BANK/server/internal/constants"
	"github.com/Eicap/EICAP-BANK/server/internal/enum"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, passOne string) {
	seedUser(db, constants.AdminOneID, "Admin One", "admin1@eicap.com", passOne, enum.Admin)
}

func seedUser(db *gorm.DB, id uuid.UUID, name, email, password string, role enum.Permission) {
	emailCopy := email
	if password == "" {
		log.Println("Skipping %S: no password configured", email)
		return
	}

	var existing model.User
	err := db.Where("id = ?", id).Take(&existing).Error
	if err == nil {
		log.Printf("Skipping %s: already exists", email)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("Error checking user %s: %v", email, err)
	}

	hashedPassword, err := pkg.Hash(password)
	if err != nil {
		log.Fatalf("Failed to hash password for %s: %v", email, err)
	}

	user := model.User{
		ID:       id,
		Name:     name,
		Email:    &emailCopy,
		Password: hashedPassword,
		Avatar:   "",
		Role:     role,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create user %s: %v", email, err)
	}
	log.Printf("User %s created successfully (%s)", email, role)
}
