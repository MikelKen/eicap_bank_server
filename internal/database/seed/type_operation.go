package seed

import (
	"errors"
	"log"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"gorm.io/gorm"
)

func SeedTypeOperations(db *gorm.DB) {
	// Conserva las operaciones históricas cambiando solo su catálogo asociado.
	db.Model(&model.TypeOperation{}).Where("code = ?", "ING").Updates(map[string]any{
		"code": "DEPO", "name": "Depósito", "cash_flow_type": "ING",
	})
	db.Model(&model.TypeOperation{}).Where("code = ?", "EGR").Updates(map[string]any{
		"code": "RETI", "name": "Retiro", "cash_flow_type": "EGR",
	})

	seedTypeOperation(db, "DEPO", "Depósito", "ING")
	seedTypeOperation(db, "RETI", "Retiro", "EGR")
	seedTypeOperation(db, "APC", "Apertura de Cuenta", "")
	seedTypeOperation(db, "APCA", "Apertura de Caja", "")
	seedTypeOperation(db, "CICA", "Cierre de Caja", "")
}

func seedTypeOperation(db *gorm.DB, code, name, cashFlowType string) {
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
		Code:         code,
		Name:         name,
		CashFlowType: cashFlowType,
	}

	if err := db.Create(&typeOperation).Error; err != nil {
		log.Fatalf("Failed to create type operation %s: %v", code, err)
	}
	log.Printf("Type operation %s (%s) created successfully", code, name)
}
