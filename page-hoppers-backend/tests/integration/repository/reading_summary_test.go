package integration_respository_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"time"
	"fmt"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"

    "page-hoppers-backend/internal/handlers"
	"page-hoppers-backend/internal/models"
	"page-hoppers-backend/tests"
)

func ParseReadingSummary(summary map[string]interface{}) (
	currentBook map[string]interface{},
	lastCompletedBook map[string]interface{},
	booksCompletedThisMonth int,
	booksCompletedThisYear int,
	totalCompletedBooks int,
	totalUncompletedBooks int,
) {
	if cb, ok := summary["currentBook"].(map[string]interface{}); ok {
		currentBook = cb
	}

	if lcb, ok := summary["lastCompletedBook"].(map[string]interface{}); ok {
		lastCompletedBook = lcb
	}

	if bctm, ok := summary["booksCompletedThisMonth"].(float64); ok {
		booksCompletedThisMonth = int(bctm)
	}

	if bcty, ok := summary["booksCompletedThisYear"].(float64); ok {
		booksCompletedThisYear = int(bcty)
	}

	if tcb, ok := summary["totalCompletedBooks"].(float64); ok {
		totalCompletedBooks = int(tcb)
	}

	if tub, ok := summary["totalUncompletedBooks"].(float64); ok {
		totalUncompletedBooks = int(tub)
	}

	return
}

// TestGetReadingSummaryIntegration tests the full flow from route → handler → database → repository
func TestGetReadingSummaryIntegration(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()

    db := tests.SetupTestDB()

	// Register route
    router.GET("/children/:id/summary", handlers.GetReadingSummaryHandler(db))

    // Seed test data
	parent := tests.CreateTestParent(db, "Alice", "alice@example.com", "password123")
    child := tests.CreateTestChild(db, "Sophie", 10, parent.ID, "1234")
    tests.SeedTestBooks(db, child.ID)

    req, _ := http.NewRequest("GET", fmt.Sprintf("/children/%d/summary", child.ID), nil)
    resp := httptest.NewRecorder()

    router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)

    var summary map[string]interface{}
    err := json.Unmarshal(resp.Body.Bytes(), &summary)
    assert.NoError(t, err)

	currentBook, lastCompletedBook, booksCompletedThisMonth, booksCompletedThisYear, totalCompletedBooks, totalUncompletedBooks := ParseReadingSummary(summary)
	// Structure checks
    assert.Contains(t, summary, "currentBook")
    assert.Contains(t, summary, "lastCompletedBook")
    assert.Contains(t, summary, "booksCompletedThisMonth")
    assert.Contains(t, summary, "booksCompletedThisYear")
    assert.Contains(t, summary, "totalCompletedBooks")
	assert.Contains(t, summary, "totalUncompletedBooks")

	// ✅ Value checks (adjust these to match your seeded data)
	assert.Equal(t, "Matilda", currentBook["title"])
	assert.Equal(t, "The Worst Witch", lastCompletedBook["title"])
	assert.Equal(t, 4, booksCompletedThisMonth)
	assert.Equal(t, 8, booksCompletedThisYear)
	assert.Equal(t, 9, totalCompletedBooks)
	assert.Equal(t, 1, totalUncompletedBooks)
}

func TestGetReadingSummary_OneStartedBook(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()

    db := tests.SetupTestDB()

    // Create parent and child
    parent := tests.CreateTestParent(db, "Bob", "bob@example.com", "password123")
    child := tests.CreateTestChild(db, "Charlie", 8, parent.ID, "5678")

    // Seed only one started book
    db.Create(&models.ReadingLog{
        ChildID: child.ID,
        Title:   "Charlie and the Chocolate Factory",
        Author:  "Roald Dahl",
        Status:  "started",
        Date:    time.Now(),
    })

    // Attach the handler
    router.GET("/children/:id/summary", handlers.GetReadingSummaryHandler(db))

    req, _ := http.NewRequest("GET", fmt.Sprintf("/children/%d/summary", child.ID), nil)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)

    var summary map[string]interface{}
    err := json.Unmarshal(resp.Body.Bytes(), &summary)
    assert.NoError(t, err)

	currentBook, lastCompletedBook, booksCompletedThisMonth, booksCompletedThisYear, totalCompletedBooks, totalUncompletedBooks := ParseReadingSummary(summary)

    // Check the fields
    assert.Equal(t, "Charlie and the Chocolate Factory", currentBook["title"])
    // Since no completed books:
    assert.Nil(t, lastCompletedBook)
    assert.Equal(t, 0, booksCompletedThisMonth)
    assert.Equal(t, 0, booksCompletedThisYear)
    assert.Equal(t, 0, totalCompletedBooks)
    assert.Equal(t, 1, totalUncompletedBooks)
}