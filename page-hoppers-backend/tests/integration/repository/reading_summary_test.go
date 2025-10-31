package integration_respository_test

import (
	"net/http"
	"testing"
	"time"

	"page-hoppers-backend/internal/models"
	"page-hoppers-backend/tests"

	"github.com/stretchr/testify/assert"
)

func ParseReadingSummary(summary map[string]interface{}) (
	currentBook map[string]interface{},
	lastCompletedBook map[string]interface{},
	totalBooksReadThisMonth int,
	totalBooksReadThisYear int,
	totalCompletedBooks int,
	totalUncompletedBooks int,
) {
	if cb, ok := summary["currentBook"].(map[string]interface{}); ok {
		currentBook = cb
	}

	if lcb, ok := summary["lastCompletedBook"].(map[string]interface{}); ok {
		lastCompletedBook = lcb
	}

	if bctm, ok := summary["totalBooksReadThisMonth"].(float64); ok {
		totalBooksReadThisMonth = int(bctm)
	}

	if bcty, ok := summary["totalBooksReadThisYear"].(float64); ok {
		totalBooksReadThisYear = int(bcty)
	}

	if tcb, ok := summary["totalCompletedBooks"].(float64); ok {
		totalCompletedBooks = int(tcb)
	}

	if tub, ok := summary["totalUncompletedBooks"].(float64); ok {
		totalUncompletedBooks = int(tub)
	}

	return
}

// No Books started or complated by the child
func TestGetReadingSummary_NoBooks(t *testing.T) {
	setup := tests.SetupSummaryTest()

	resp := setup.GetSummary()
	assert.Equal(t, http.StatusOK, resp.Code)

	summary := tests.ParseSummary(t, resp)
	currentBook, lastCompletedBook, totalBooksReadThisMonth,
		totalBooksReadThisYear, totalCompletedBooks, totalUncompletedBooks :=
		ParseReadingSummary(summary)

	// Validate results
	assert.Equal(t, 0, totalBooksReadThisMonth, "no books read this month")
	assert.Equal(t, 0, totalBooksReadThisYear, "no books read this year")
	assert.Equal(t, 0, totalCompletedBooks, "no total completed books")
	assert.Equal(t, 0, totalUncompletedBooks, "no total uncompleted books")
	assert.Nil(t, currentBook, "expected no current book")
	assert.Nil(t, lastCompletedBook, "expected no last completed book")
}

// Only one book has been started by the child.
func TestGetReadingSummary_OneStartedBook(t *testing.T) {
	setup := tests.SetupSummaryTest()

	setup.DB.Create(&models.ReadingLog{
		ChildID: setup.Child.ID,
		Title:   "Matilda",
		Author:  "Roald Dahl",
		Status:  "started",
		Date:    time.Now(),
	})

	resp := setup.GetSummary()
	assert.Equal(t, http.StatusOK, resp.Code)

	summary := tests.ParseSummary(t, resp)
	currentBook, lastCompletedBook, totalBooksReadThisMonth,
		totalBooksReadThisYear, totalCompletedBooks, totalUncompletedBooks :=
		ParseReadingSummary(summary)

	// ✅ Assertions
	assert.NotNil(t, currentBook, "expected current book to be populated")
	assert.Nil(t, lastCompletedBook, "expected no last completed book")
	assert.Equal(t, 0, totalBooksReadThisMonth, "no books completed this month")
	assert.Equal(t, 0, totalBooksReadThisYear, "no books completed this year")
	assert.Equal(t, 0, totalCompletedBooks, "no total completed books")
	assert.Equal(t, 1, totalUncompletedBooks, "one uncompleted book total")
}

// Only one book has been completed by the child.
func TestGetReadingSummary_OneCompletedBook(t *testing.T) {
	setup := tests.SetupSummaryTest()

	setup.DB.Create(&models.ReadingLog{
		ChildID: setup.Child.ID,
		Title:   "The Railway Children",
		Author:  "E. B. Nesbit",
		Status:  "completed",
		Date:    time.Now(),
	})

	resp := setup.GetSummary()
	assert.Equal(t, http.StatusOK, resp.Code)

	summary := tests.ParseSummary(t, resp)
	currentBook, lastCompletedBook, totalBooksReadThisMonth,
		totalBooksReadThisYear, totalCompletedBooks, totalUncompletedBooks :=
		ParseReadingSummary(summary)

	// ✅ Assertions
	assert.Nil(t, currentBook, "expected no current book")
	assert.NotNil(t, lastCompletedBook, "expected a last completes book")
	assert.Equal(t, 1, totalBooksReadThisMonth, "no books completed this month")
	assert.Equal(t, 1, totalBooksReadThisYear, "no books completed this year")
	assert.Equal(t, 1, totalCompletedBooks, "no total completed books")
	assert.Equal(t, 0, totalUncompletedBooks, "one uncompleted book total")
}

// One Book completed last month and one book completed this mmonth.
func TestGetReadingSummary_CompletedBooksAcrossMonths(t *testing.T) {
	setup := tests.SetupSummaryTest()

	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastMonth := thisMonth.AddDate(0, -1, 0) // first day of previous month

	// 🪄 Add one book completed THIS month
	setup.DB.Create(&models.ReadingLog{
		ChildID: setup.Child.ID,
		Title:   "Matilda",
		Author:  "Roald Dahl",
		Status:  "completed",
		Date:    thisMonth,
	})

	// 🪄 Add one book completed LAST month
	setup.DB.Create(&models.ReadingLog{
		ChildID: setup.Child.ID,
		Title:   "The BFG",
		Author:  "Roald Dahl",
		Status:  "completed",
		Date:    lastMonth,
	})

	resp := setup.GetSummary()
	assert.Equal(t, http.StatusOK, resp.Code)

	summary := tests.ParseSummary(t, resp)
	currentBook, lastCompletedBook, totalBooksReadThisMonth,
		totalBooksReadThisYear, totalCompletedBooks, totalUncompletedBooks :=
		ParseReadingSummary(summary)

	// ✅ Assertions
	assert.Nil(t, currentBook, "no current book should be set (all completed)")
	assert.NotNil(t, lastCompletedBook, "expected a last completed book")

	assert.Equal(t, 1, totalBooksReadThisMonth, "expected one book completed this month")
	// totalBooksReadThisYear depends on whether lastMonth is in the same year
	expectedBooksThisYear := 1
	if lastMonth.Year() == thisMonth.Year() {
		expectedBooksThisYear = 2
	}
	assert.Equal(t, expectedBooksThisYear, totalBooksReadThisYear, "total books completed this year")

	assert.Equal(t, 2, totalCompletedBooks, "expected total of two completed books")
	assert.Equal(t, 0, totalUncompletedBooks, "expected no uncompleted books")

	assert.Equal(t, "Matilda", lastCompletedBook["title"], "expected Matilda to be last completed")
}

func TestGetReadingSummary_EightBooksThreeMonthsWithCurrent(t *testing.T) {
	setup := tests.SetupSummaryTest()

	today := time.Now()
	currentYear := today.Year()
	currentMonth := today.Month()

	// Calculate months
	firstOfThisMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	lastMonth := firstOfThisMonth.AddDate(0, -1, 0)
	twoMonthsAgo := firstOfThisMonth.AddDate(0, -2, 0)

	// 🪄 All books in one slice
	books := []models.ReadingLog{
		// Current book in progress
		{ChildID: setup.Child.ID, Title: "Charlotte's Web", Author: "E. B. White", Status: "started", Date: today},

		// 3 completed books this month
		{ChildID: setup.Child.ID, Title: "Matilda", Author: "Roald Dahl", Status: "completed", Date: today},
		{ChildID: setup.Child.ID, Title: "The Worst Witch", Author: "Jill Murphy", Status: "completed", Date: today.AddDate(0, 0, -1)},
		{ChildID: setup.Child.ID, Title: "The Tale of Peter Rabbit", Author: "Beatrix Potter", Status: "completed", Date: today.AddDate(0, 0, -2)},

		// 3 completed books last month
		{ChildID: setup.Child.ID, Title: "The BFG", Author: "Roald Dahl", Status: "completed", Date: lastMonth},
		{ChildID: setup.Child.ID, Title: "Charlotte's Web", Author: "E. B. White", Status: "completed", Date: lastMonth.AddDate(0, 0, 1)},
		{ChildID: setup.Child.ID, Title: "Fantastic Mr Fox", Author: "Roald Dahl", Status: "completed", Date: lastMonth.AddDate(0, 0, 2)},

		// 2 completed books two months ago
		{ChildID: setup.Child.ID, Title: "The Lion, the Witch and the Wardrobe", Author: "C. S. Lewis", Status: "completed", Date: twoMonthsAgo},
		{ChildID: setup.Child.ID, Title: "The Secret Garden", Author: "Frances Hodgson Burnett", Status: "completed", Date: twoMonthsAgo.AddDate(0, 0, 1)},
	}

	// 🪄 Insert all books in one go
	setup.DB.Create(&books)

	// 🧭 GET summary
	resp := setup.GetSummary()
	assert.Equal(t, http.StatusOK, resp.Code)

	summary := tests.ParseSummary(t, resp)
	currentBook, lastCompletedBook, totalBooksReadThisMonth,
		totalBooksReadThisYear, totalCompletedBooks, totalUncompletedBooks :=
		ParseReadingSummary(summary)

	// ✅ Assertions
	assert.NotNil(t, currentBook, "expected a current book")
	assert.Equal(t, "Charlotte's Web", currentBook["title"], "current book should be correct")
	assert.NotNil(t, lastCompletedBook, "expected a last completed book")
	assert.Equal(t, "Matilda", lastCompletedBook["title"], "last completed book should be most recent")

	assert.Equal(t, 3, totalBooksReadThisMonth, "3 books completed this month")

	// Calculate expected books this year
	expectedThisYear := 0
	for _, logDate := range []time.Time{
		today, today.AddDate(0, 0, -1), today.AddDate(0, 0, -2), // this month
		lastMonth, lastMonth.AddDate(0, 0, 1), lastMonth.AddDate(0, 0, 2), // last month
		twoMonthsAgo, twoMonthsAgo.AddDate(0, 0, 1), // two months ago
	} {
		if logDate.Year() == currentYear {
			expectedThisYear++
		}
	}
	assert.Equal(t, expectedThisYear, totalBooksReadThisYear, "total books completed this year")
	assert.Equal(t, 8, totalCompletedBooks, "total completed books")
	assert.Equal(t, 1, totalUncompletedBooks, "one uncompleted book")
}

func TestGetReadingSummary_TwoStartedBooks(t *testing.T) {
	setup := tests.SetupSummaryTest()

	now := time.Now()

	// Add two started books
	setup.DB.Create(&models.ReadingLog{
		ChildID: setup.Child.ID,
		Title:   "Charlotte's Web",
		Author:  "E. B. White",
		Status:  "started",
		Date:    now.AddDate(0, 0, -2), // started 2 days ago
	})

	setup.DB.Create(&models.ReadingLog{
		ChildID: setup.Child.ID,
		Title:   "Matilda",
		Author:  "Roald Dahl",
		Status:  "started",
		Date:    now, // started today
	})

	resp := setup.GetSummary()
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse response
	summary := tests.ParseSummary(t, resp)
	currentBook, lastCompletedBook, totalBooksReadThisMonth,
		totalBooksReadThisYear, totalCompletedBooks, totalUncompletedBooks :=
		ParseReadingSummary(summary)

	// Assertions
	assert.NotNil(t, currentBook, "expected current book to be populated")
	assert.Nil(t, lastCompletedBook, "expected no last completed book")

	// The latest started book should be returned
	assert.Equal(t, "Matilda", currentBook["title"], "expected Matilda to be the current book")

	assert.Equal(t, 0, totalBooksReadThisMonth)
	assert.Equal(t, 0, totalBooksReadThisYear)
	assert.Equal(t, 0, totalCompletedBooks)
	assert.Equal(t, 2, totalUncompletedBooks)
}
