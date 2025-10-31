package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"page-hoppers-backend/internal/handlers"
	"page-hoppers-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

type SummaryTestSetup struct {
	DB     *gorm.DB
	Router *gin.Engine
	Parent *models.User
	Child  *models.User
}

// SetupSummaryTest sets up a test DB, parent, child, and router with /children/:id/summary route
func SetupSummaryTest() *SummaryTestSetup {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db := SetupTestDB()
	parent := CreateTestParent(db, "Bob", "bob@example.com", "password123")
	child := CreateTestChild(db, "Charlie", 8, parent.ID, "5678")

	handler := handlers.ReadingLogHandler{DB: db}
	router.GET("/children/:id/summary", handler.GetReadingSummary)

	return &SummaryTestSetup{
		DB:     db,
		Router: router,
		Parent: parent,
		Child:  child,
	}
}

// GetSummary sends a GET request to /children/:id/summary and returns the response
func (s *SummaryTestSetup) GetSummary() *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/children/%d/summary", s.Child.ID), nil)
	resp := httptest.NewRecorder()
	s.Router.ServeHTTP(resp, req)
	return resp
}

// ParseSummary parses the response into a generic map for flexible access
func ParseSummary(t *testing.T, resp *httptest.ResponseRecorder) map[string]interface{} {
	var summary map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &summary)
	assert.NoError(t, err, "failed to parse reading summary JSON")
	return summary
}
