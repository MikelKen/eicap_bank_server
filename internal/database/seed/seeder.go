package seed

import "gorm.io/gorm"

func Run(db *gorm.DB, adminPassOne string) {
	SeedUsers(db, adminPassOne)
	SeedDenominations(db)
	SeedTypeAccounts(db)
	SeedTypeOperations(db)
}
