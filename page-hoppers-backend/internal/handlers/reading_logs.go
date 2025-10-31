package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"page-hoppers-backend/internal/models"
)

type ReadingLogHandler struct {
	DB *gorm.DB
}

func NewReadingLogHandler(db *gorm.DB) *ReadingLogHandler {
	return &ReadingLogHandler{
		DB: db,
	}
}

// ---------------------------
// Request/Response structs
type CreateReadingLogRequest struct {
	Title          string `json:"title"`
	Author         string `json:"author,omitempty"`
	Status         string `json:"status"` // "started" or "completed"
	Date           string `json:"date"`   // ISO date string
	OpenLibraryKey string `json:"open_library_key,omitempty"`
	CoverID        *int   `json:"cover_id,omitempty"`
}

type ReadingLogResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Author         string    `json:"author,omitempty"`
	Status         string    `json:"status"`
	Date           time.Time `json:"date"`
	OpenLibraryKey string    `json:"open_library_key,omitempty"`
	CoverID        *int      `json:"cover_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ---------------------------
// Create a reading log (child)
func (h *ReadingLogHandler) CreateReadingLog(c *gin.Context) {
	childIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	childID := childIDValue.(uint)

	var req CreateReadingLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Title == "" || req.Status == "" || req.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, status, and date are required"})
		return
	}

	if req.Status != "started" && req.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'started' or 'completed'"})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	var child models.User
	if err := h.DB.Where("id = ? AND role = ?", childID, "child").First(&child).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child not found"})
		return
	}

	readingLog := models.ReadingLog{
		Title:          req.Title,
		Author:         req.Author,
		Status:         req.Status,
		Date:           date,
		ChildID:        childID,
		OpenLibraryKey: req.OpenLibraryKey,
		CoverID:        req.CoverID,
	}

	if err := h.DB.Create(&readingLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reading log"})
		return
	}

	c.JSON(http.StatusOK, ReadingLogResponse{
		ID:             readingLog.ID,
		Title:          readingLog.Title,
		Author:         readingLog.Author,
		Status:         readingLog.Status,
		Date:           readingLog.Date,
		OpenLibraryKey: readingLog.OpenLibraryKey,
		CoverID:        readingLog.CoverID,
		CreatedAt:      readingLog.CreatedAt,
	})
}

// ---------------------------
// Get all reading logs for a child
func (h *ReadingLogHandler) GetReadingLogs(c *gin.Context) {
	childIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	childID := childIDValue.(uint)

	var child models.User
	if err := h.DB.Where("id = ? AND role = ?", childID, "child").First(&child).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child not found"})
		return
	}

	var logs []models.ReadingLog
	if err := h.DB.Where("child_id = ?", childID).Order("date DESC, created_at DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reading logs"})
		return
	}

	var responses []ReadingLogResponse
	for _, log := range logs {
		responses = append(responses, ReadingLogResponse{
			ID:             log.ID,
			Title:          log.Title,
			Author:         log.Author,
			Status:         log.Status,
			Date:           log.Date,
			OpenLibraryKey: log.OpenLibraryKey,
			CoverID:        log.CoverID,
			CreatedAt:      log.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

// ---------------------------
// Get reading logs for a specific child (parent access)
func (h *ReadingLogHandler) GetChildReadingLogs(c *gin.Context) {
	parentIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	parentID := parentIDValue.(uint)

	childIDStr := c.Query("child_id")
	if childIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Child ID is required"})
		return
	}

	childIDUint, err := strconv.ParseUint(childIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid child ID"})
		return
	}
	childID := uint(childIDUint)

	var child models.User
	if err := h.DB.Where("id = ? AND parent_id = ? AND role = ?", childID, parentID, "child").First(&child).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child not found or unauthorized"})
		return
	}

	var logs []models.ReadingLog
	if err := h.DB.Where("child_id = ?", childID).Order("date DESC, created_at DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reading logs"})
		return
	}

	var responses []ReadingLogResponse
	for _, log := range logs {
		responses = append(responses, ReadingLogResponse{
			ID:             log.ID,
			Title:          log.Title,
			Author:         log.Author,
			Status:         log.Status,
			Date:           log.Date,
			OpenLibraryKey: log.OpenLibraryKey,
			CoverID:        log.CoverID,
			CreatedAt:      log.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}