package repository

import (
	"time"
	"gorm.io/gorm"
	"page-hoppers-backend/internal/models"
)

type ReadingSummary struct {
	CurrentBook    *models.ReadingLog `json:"currentBook,omitempty"`
	LastCompletedBook       *models.ReadingLog `json:"lastCompletedBook,omitempty"`
	TotalUncompletedBooks	int					`json:"totalUncompletedBooks"`
	BooksCompletedThisMonth int                `json:"booksCompletedThisMonth"`
	BooksCompletedThisYear  int                `json:"booksCompletedThisYear"`
	TotalCompletedBooks     int                `json:"totalCompletedBooks"`
}

// db *gorm.DB → a pointer to the GORM database connection.
// childID uint → the unique ID of the child whose reading summary we’re fetching.
func GetReadingSummary(db *gorm.DB, childID uint) (*ReadingSummary, error) {
	var logs []models.ReadingLog
	if err := db.Where("child_id = ?", childID).Order("date desc").Find(&logs).Error; err != nil {
		return nil, err
	}

	summary := &ReadingSummary{}
	now := time.Now()
	currentMonth := now.Month()
	currentYear := now.Year()

	for _, log := range logs {
		if log.Status == "started" {
			summary.TotalUncompletedBooks++
		}
		if summary.CurrentBook == nil && log.Status == "started" {
			summary.CurrentBook = &log
		}
		if summary.LastCompletedBook == nil && log.Status == "completed" {
			summary.LastCompletedBook = &log
		}
		if log.Status == "completed" {
			summary.TotalCompletedBooks++
			if log.Date.Month() == currentMonth && log.Date.Year() == currentYear {
				summary.BooksCompletedThisMonth++
			}
			if log.Date.Year() == currentYear {
				summary.BooksCompletedThisYear++
			}
		}
	}

	return summary, nil
}