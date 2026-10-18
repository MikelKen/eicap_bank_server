package seed

import "gorm.io/gorm"

func Run(db *gorm.DB, adminPassOne string, eicapPass string) {
	SeedUsers(db, adminPassOne, eicapPass)
	SeedDenominations(db)
	SeedTypeAccounts(db)
	SeedTypeOperations(db)
}
