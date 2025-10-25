package tests

import (
	"time"
	"fmt"
	
	"page-hoppers-backend/internal/models"
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
		Name:     name,
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

func SeedTestBooks(db *gorm.DB, childID uint) {
	now := time.Now()
    logs := []models.ReadingLog{
        // This month
		{ChildID: childID, Title: "The Worst Witch", Author: "Jill Murphy", Status: "completed", Date: now.AddDate(0, 0, -1)}, // yesterday
		{ChildID: childID, Title: "Matilda", Author: "Roald Dahl", Status: "started", Date: now}, // today

		// Earlier this month
		{ChildID: childID, Title: "Charlotte's Web", Author: "E. B. White", Status: "completed", Date: now.AddDate(0, 0, -7)},
		{ChildID: childID, Title: "The BFG", Author: "Roald Dahl", Status: "completed", Date: now.AddDate(0, 0, -10)},
		{ChildID: childID, Title: "The Lion, the Witch and the Wardrobe", Author: "C. S. Lewis", Status: "completed", Date: now.AddDate(0, 0, -14)},

		// Last month
		{ChildID: childID, Title: "Harry Potter and the Philosopher's Stone", Author: "J. K. Rowling", Status: "completed", Date: now.AddDate(0, -1, 0)},
		{ChildID: childID, Title: "Fantastic Mr Fox", Author: "Roald Dahl", Status: "completed", Date: now.AddDate(0, -1, -3)},

		// Earlier this year
		{ChildID: childID, Title: "The Secret Garden", Author: "Frances Hodgson Burnett", Status: "completed", Date: now.AddDate(0, -3, 0)},
		{ChildID: childID, Title: "The Witches", Author: "Roald Dahl", Status: "completed", Date: now.AddDate(0, -6, 0)},

		// Last Year
		{ChildID: childID, Title: "The Railway Children", Author: "E. Nesbit", Status: "completed", Date: now.AddDate(-1, -2, 0)},
    }
    for _, log := range logs {
        db.Create(&log)
    }

	var count int64
    db.Model(&models.ReadingLog{}).Where("child_id = ?", childID).Count(&count)
    fmt.Println("Seeded books count:", count)
}