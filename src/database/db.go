package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	// Setup database
	db, _ := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	return db
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Transaction{}, &Passbook{})
}
