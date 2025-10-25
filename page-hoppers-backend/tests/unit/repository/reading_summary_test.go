package unit_repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"page-hoppers-backend/internal/models"
	"page-hoppers-backend/internal/repository"
)

// setupTestDB creates an in-memory SQLite database for testing so that tests run fast and don't affect real data.
// Runs migrations so your User and ReadingLog tables exist.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.ReadingLog{})
	assert.NoError(t, err)

	return db
}

// TEST 1: TestGetReadingSummary_EmptyLogs tests when there are no reading logs
// Defines a test function in Go.
// All test functions must start with Test and take a single argument t *testing.T.
// The Go test runner (go test) looks for functions with this signature to run them automatically.
func TestGetReadingSummary_EmptyLogs(t *testing.T) {
	// Calls a helper function that creates a brand new in-memory database just for this test.
	db := setupTestDB(t)
	// Sets up a fake child ID (with a value of 1) to use when calling the function.
	childID := uint(1)

	// Calls the actual function you’re testing — GetReadingSummary.
	summary, err := repository.GetReadingSummary(db, childID)

	// Uses testify/assert to check that no error was returned.
	assert.NoError(t, err)
	// Checks that summary is not nil.
	assert.NotNil(t, summary)
	// Verify that both CurrentBook and LastBook are nil.
	assert.Nil(t, summary.CurrentBook)
	assert.Nil(t, summary.LastBook)
	// Check that all count fields in the summary are 0.
	assert.Equal(t, 0, summary.TotalBooks)
	assert.Equal(t, 0, summary.BooksThisMonth)
	assert.Equal(t, 0, summary.BooksThisYear)
}

// TEST 2: TestGetReadingSummary_WithData tests with various reading logs
func TestGetReadingSummary_WithData(t *testing.T) {
	db := setupTestDB(t)
	childID := uint(1)
	// Stores the current date/time (now) for calculating relative dates.
	// You’ll use now with AddDate() to generate logs from different times — yesterday, last month, last year, etc.
	now := time.Now()

	// Create test data
	readingLogs := []models.ReadingLog{
		// This entry simulates a current book that’s still being read (Status: "started").
		{
			ChildID: childID,
			Title:   "The Railway Children",
			Author:  "E. B. Nesbit",
			Status:  "started",
			Date:    now.AddDate(0, 0, -1), // Yesterday
		},
		{
			ChildID: childID,
			Title:   "The Worst Witch",
			Author:  "Jill Murphy",
			Status:  "completed",
			Date:    now.AddDate(0, 0, -2), // 2 days ago
		},
		{
			ChildID: childID,
			Title:   "This Month Book",
			Author:  "Author C",
			Status:  "completed",
			Date:    now.AddDate(0, 0, -5), // 5 days ago (this month)
		},
		{
			ChildID: childID,
			Title:   "This Year Book",
			Author:  "Author D",
			Status:  "completed",
			Date:    now.AddDate(0, -1, 0), // Last month (this year)
		},
		{
			ChildID: childID,
			Title:   "Last Year Book",
			Author:  "Author E",
			Status:  "completed",
			Date:    now.AddDate(-1, 0, 0), // Last year
		},
	}

	// Insert test data
	for _, log := range readingLogs {
		result := db.Create(&log)
		assert.NoError(t, result.Error)
	}

	summary, err := repository.GetReadingSummary(db, childID)

	assert.NoError(t, err)
	assert.NotNil(t, summary)

	// Check current book (most recent "started" book)
	assert.NotNil(t, summary.CurrentBook)
	assert.Equal(t, "The Railway Children", summary.CurrentBook.Title)
	assert.Equal(t, "E. B. Nesbit", summary.CurrentBook.Author)
	assert.Equal(t, "started", summary.CurrentBook.Status)

	// // Check last book (most recent "completed" book)
	// assert.NotNil(t, summary.LastBook)
	// assert.Equal(t, "Last Completed Book", summary.LastBook.Title)
	// assert.Equal(t, "completed", summary.LastBook.Status)

	// // Check counts
	// assert.Equal(t, 4, summary.TotalBooks)        // 4 completed books
	// assert.Equal(t, 2, summary.BooksThisMonth) // 2 books this month
	// assert.Equal(t, 3, summary.BooksThisYear)   // 3 books this year
}

// // TestGetReadingSummary_OnlyCompletedBooks tests when there are only completed books
// func TestGetReadingSummary_OnlyCompletedBooks(t *testing.T) {
// 	db := setupTestDB(t)
// 	childID := uint(1)
// 	now := time.Now()

// 	readingLogs := []models.ReadingLog{
// 		{
// 			ChildID: childID,
// 			Title:   "Completed Book 1",
// 			Status:  "completed",
// 			Date:    now.AddDate(0, 0, -1),
// 		},
// 		{
// 			ChildID: childID,
// 			Title:   "Completed Book 2",
// 			Status:  "completed",
// 			Date:    now.AddDate(0, 0, -2),
// 		},
// 	}

// 	for _, log := range readingLogs {
// 		result := db.Create(&log)
// 		assert.NoError(t, result.Error)
// 	}

// 	summary, err := repository.GetReadingSummary(db, childID)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, summary)
// 	assert.Nil(t, summary.CurrentBook) // No started books
// 	assert.NotNil(t, summary.LastBook)
// 	assert.Equal(t, "Completed Book 1", summary.LastBook.Title)
// 	assert.Equal(t, 2, summary.TotalBooks)
// }

// // TestGetReadingSummary_OnlyStartedBooks tests when there are only started books
// func TestGetReadingSummary_OnlyStartedBooks(t *testing.T) {
// 	db := setupTestDB(t)
// 	childID := uint(1)
// 	now := time.Now()

// 	readingLogs := []models.ReadingLog{
// 		{
// 			ChildID: childID,
// 			Title:   "Started Book 1",
// 			Status:  "started",
// 			Date:    now.AddDate(0, 0, -1),
// 		},
// 		{
// 			ChildID: childID,
// 			Title:   "Started Book 2",
// 			Status:  "started",
// 			Date:    now.AddDate(0, 0, -2),
// 		},
// 	}

// 	for _, log := range readingLogs {
// 		result := db.Create(&log)
// 		assert.NoError(t, result.Error)
// 	}

// 	summary, err := repository.GetReadingSummary(db, childID)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, summary)
// 	assert.NotNil(t, summary.CurrentBook) // Most recent started book
// 	assert.Equal(t, "Started Book 1", summary.CurrentBook.Title)
// 	assert.Nil(t, summary.LastBook) // No completed books
// 	assert.Equal(t, 0, summary.TotalBooks)
// }

// // TestGetReadingSummary_MultipleChildren tests that it only returns data for the specified child
// func TestGetReadingSummary_MultipleChildren(t *testing.T) {
// 	db := setupTestDB(t)
// 	childID1 := uint(1)
// 	childID2 := uint(2)
// 	now := time.Now()

// 	// Create logs for child 1
// 	logs1 := []models.ReadingLog{
// 		{
// 			ChildID: childID1,
// 			Title:   "Child 1 Book",
// 			Status:  "completed",
// 			Date:    now.AddDate(0, 0, -1),
// 		},
// 	}

// 	// Create logs for child 2
// 	logs2 := []models.ReadingLog{
// 		{
// 			ChildID: childID2,
// 			Title:   "Child 2 Book",
// 			Status:  "completed",
// 			Date:    now.AddDate(0, 0, -1),
// 		},
// 	}

// 	allLogs := append(logs1, logs2...)
// 	for _, log := range allLogs {
// 		result := db.Create(&log)
// 		assert.NoError(t, result.Error)
// 	}

// 	// Test child 1 summary
// 	summary1, err := repository.GetReadingSummary(db, childID1)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, summary1)
// 	assert.Equal(t, 1, summary1.TotalBooks)
// 	assert.Equal(t, "Child 1 Book", summary1.LastBook.Title)

// 	// Test child 2 summary
// 	summary2, err := repository.GetReadingSummary(db, childID2)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, summary2)
// 	assert.Equal(t, 1, summary2.TotalBooks)
// 	assert.Equal(t, "Child 2 Book", summary2.LastBook.Title)
// }