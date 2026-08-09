package seed

import (
	"errors"
	"log"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func SeedDenominations(db *gorm.DB) {
	seedDenomination(db, "bill", 200, "Billete de 200")
	seedDenomination(db, "bill", 100, "Billete de 100")
	seedDenomination(db, "bill", 50, "Billete de 50")
	seedDenomination(db, "bill", 20, "Billete de 20")
	seedDenomination(db, "bill", 10, "Billete de 10")
	seedDenomination(db, "coin", 5, "Moneda de 5")
	seedDenomination(db, "coin", 2, "Moneda de 2")
	seedDenomination(db, "coin", 1, "Moneda de 1")
}

func seedDenomination(db *gorm.DB, typeDenom string, value int64, name string) {
	var existing model.Denomination
	err := db.Where("type = ? AND value = ?", typeDenom, value).Take(&existing).Error
	if err == nil {
		log.Printf("Skipping denomination %s %s: already exists", typeDenom, name)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("Error checking denomination %s %s: %v", typeDenom, name, err)
	}

	denomination := model.Denomination{
		Type:  typeDenom,
		Value: decimal.NewFromInt(value),
		Name:  name,
	}

	if err := db.Create(&denomination).Error; err != nil {
		log.Fatalf("Failed to create denomination %s %s: %v", typeDenom, name, err)
	}
	log.Printf("Denomination %s %s created successfully", typeDenom, name)
}
