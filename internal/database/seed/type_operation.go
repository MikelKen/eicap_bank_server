package seed

import (
	"errors"
	"log"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"gorm.io/gorm"
)

func SeedTypeOperations(db *gorm.DB) {
	seedTypeOperation(db, "ING", "Ingreso")
	seedTypeOperation(db, "EGR", "Egreso")
	seedTypeOperation(db, "APC", "Apertura de Cuenta")
	seedTypeOperation(db, "APCA", "Apertura de Caja")
	seedTypeOperation(db, "CICA", "Cierre de Caja")
}

func seedTypeOperation(db *gorm.DB, code, name string) {
	var existing model.TypeOperation
	err := db.Where("code = ?", code).Take(&existing).Error
	if err == nil {
		log.Printf("Skipping type operation %s: already exists", code)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("Error checking type operation %s: %v", code, err)
	}

	typeOperation := model.TypeOperation{
		Code: code,
		Name: name,
	}

	if err := db.Create(&typeOperation).Error; err != nil {
		log.Fatalf("Failed to create type operation %s: %v", code, err)
	}
	log.Printf("Type operation %s (%s) created successfully", code, name)
}
