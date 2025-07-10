package handlers

import (
	"github.com/lcleigh/page-hoppers-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	
	// Auto migrate the schema
	if err := db.AutoMigrate(&models.User{}, &models.ReadingLog{}); err != nil {
		panic("failed to migrate test database")
	}
	
	return db
}

// hashPassword creates a bcrypt hash for testing
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to hash password")
	}
	return string(hashedPassword)
}

// createTestParent creates a test parent user in the database
func createTestParent(db *gorm.DB, name, email, password string) *models.User {
	parent := &models.User{
		Username: name,
		Email:    email,
		Password: hashPassword(password),
		Role:     "parent",
	}
	db.Create(parent)
	return parent
}

// createTestChild creates a test child user in the database
func createTestChild(db *gorm.DB, username, name string, age int, parentID uint, pin string) *models.User {
	child := &models.User{
		Username: username,
		Name:     name,
		Age:      age,
		PIN:      hashPassword(pin),
		Role:     "child",
		ParentID: &parentID,
	}
	db.Create(child)
	return child
} 