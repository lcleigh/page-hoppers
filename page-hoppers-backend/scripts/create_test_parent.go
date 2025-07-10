package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/lcleigh/page-hoppers-backend/models"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("DATABASE_URL environment variable not set")
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	password := "testpassword"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to hash password")
	}

	parent := models.User{
		Username: "testparent",
		Email:    "parent@example.com",
		Password: string(hashed),
		Role:     "parent",
	}

	if err := db.Where(models.User{Email: parent.Email}).FirstOrCreate(&parent).Error; err != nil {
		panic(err)
	}

	fmt.Println("Test parent created: parent@example.com / testpassword")
} 