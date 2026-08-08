package seed

import (
	"errors"
	"log"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"gorm.io/gorm"
)

func SeedTypeAccounts(db *gorm.DB) {
	seedTypeAccount(db, "Caja de Ahorro")
	seedTypeAccount(db, "Cuenta Corriente")
	seedTypeAccount(db, "Caja DPF")
}

func seedTypeAccount(db *gorm.DB, name string) {
	var existing model.TypeAccount
	err := db.Where("name = ?", name).Take(&existing).Error
	if err == nil {
		log.Printf("Skipping type account %s: already exists", name)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("Error checking type account %s: %v", name, err)
	}

	typeAccount := model.TypeAccount{
		Name: name,
	}

	if err := db.Create(&typeAccount).Error; err != nil {
		log.Fatalf("Failed to create type account %s: %v", name, err)
	}
	log.Printf("Type account %s created successfully", name)
}
