package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/lcleigh/page-hoppers-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Connect to database
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("DATABASE_URL environment variable not set")
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Auto migrate the schema
	fmt.Println("Migrating database...")
	if err := db.AutoMigrate(&models.User{}, &models.ReadingLog{}); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database migration completed successfully!")
	fmt.Println("Tables created:")
	fmt.Println("- users")
	fmt.Println("- reading_logs")
} 