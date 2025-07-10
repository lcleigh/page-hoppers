package tests

import (
	"github.com/lcleigh/page-hoppers-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() *gorm.DB {
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

// CreateTestParent creates a test parent user in the database
func CreateTestParent(db *gorm.DB, name, email, password string) *models.User {
	parent := &models.User{
		Username: name,
		Email:    email,
		Password: password, // Note: In real tests, this should be hashed
		Role:     "parent",
	}
	
	db.Create(parent)
	return parent
}

// CreateTestChild creates a test child user in the database
func CreateTestChild(db *gorm.DB, name string, age int, parentID uint, pin string) *models.User {
	child := &models.User{
		Name:     name,
		Age:      age,
		PIN:      pin, // Note: In real tests, this should be hashed
		Role:     "child",
		ParentID: &parentID,
	}
	
	db.Create(child)
	return child
} 