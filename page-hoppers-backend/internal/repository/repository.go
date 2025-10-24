package repository

import (
	"log"
	"os"

	"github.com/lcleigh/page-hoppers-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&models.User{}, &models.ReadingLog{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}