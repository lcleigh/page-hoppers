package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"github.com/lcleigh/page-hoppers-backend/internal/models"
	"gorm.io/gorm"
)

type ReadingLogHandler struct {
	db *gorm.DB
}

func NewReadingLogHandler(db *gorm.DB) *ReadingLogHandler {
	return &ReadingLogHandler{
		db: db,
	}
}

type CreateReadingLogRequest struct {
	Title           string `json:"title"`
	Author          string `json:"author,omitempty"`
	Status          string `json:"status"` // "started" or "completed"
	Date            string `json:"date"`   // ISO date string
	OpenLibraryKey  string `json:"open_library_key,omitempty"`
	CoverID         *int   `json:"cover_id,omitempty"`
}

type ReadingLogResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author,omitempty"`
	Status          string    `json:"status"`
	Date            time.Time `json:"date"`
	OpenLibraryKey  string    `json:"open_library_key,omitempty"`
	CoverID         *int      `json:"cover_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateReadingLog creates a new reading log entry for a child
func (h *ReadingLogHandler) CreateReadingLog(w http.ResponseWriter, r *http.Request) {
	// Extract child ID from JWT token
	childID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateReadingLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" || req.Status == "" || req.Date == "" {
		http.Error(w, "Title, status, and date username", http.StatusBadRequest)
		return
	}

	// Validate status
	if req.Status != "started" && req.Status != "completed" {
		http.Error(w, "Status must be 'started' or 'completed'", http.StatusBadRequest)
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Verify child exists
	var child models.User
	if err := h.db.Where("id = ? AND role = ?", childID, "child").First(&child).Error; err != nil {
		http.Error(w, "Child not found", http.StatusNotFound)
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

	if err := h.db.Create(&readingLog).Error; err != nil {
		http.Error(w, "Failed to create reading log", http.StatusInternalServerError)
		return
	}

	response := ReadingLogResponse{
		ID:             readingLog.ID,
		Title:          readingLog.Title,
		Author:         readingLog.Author,
		Status:         readingLog.Status,
		Date:           readingLog.Date,
		OpenLibraryKey: readingLog.OpenLibraryKey,
		CoverID:        readingLog.CoverID,
		CreatedAt:      readingLog.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetReadingLogs fetches all reading logs for a child
func (h *ReadingLogHandler) GetReadingLogs(w http.ResponseWriter, r *http.Request) {
	// Extract child ID from JWT token
	childID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify child exists
	var child models.User
	if err := h.db.Where("id = ? AND role = ?", childID, "child").First(&child).Error; err != nil {
		http.Error(w, "Child not found", http.StatusNotFound)
		return
	}

	var readingLogs []models.ReadingLog
	if err := h.db.Where("child_id = ?", childID).Order("date DESC, created_at DESC").Find(&readingLogs).Error; err != nil {
		http.Error(w, "Failed to fetch reading logs", http.StatusInternalServerError)
		return
	}

	var responses []ReadingLogResponse
	for _, log := range readingLogs {
		response := ReadingLogResponse{
			ID:             log.ID,
			Title:          log.Title,
			Author:         log.Author,
			Status:         log.Status,
			Date:           log.Date,
			OpenLibraryKey: log.OpenLibraryKey,
			CoverID:        log.CoverID,
			CreatedAt:      log.CreatedAt,
		}
		responses = append(responses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// GetChildReadingLogs fetches reading logs for a specific child (parent access)
func (h *ReadingLogHandler) GetChildReadingLogs(w http.ResponseWriter, r *http.Request) {
	// Extract parent ID from JWT token
	parentID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get child ID from query parameter
	childIDStr := r.URL.Query().Get("child_id")
	if childIDStr == "" {
		http.Error(w, "Child ID is required", http.StatusBadRequest)
		return
	}

	childID, err := strconv.ParseUint(childIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid child ID", http.StatusBadRequest)
		return
	}

	// Verify the child belongs to the parent
	var child models.User
	if err := h.db.Where("id = ? AND parent_id = ? AND role = ?", childID, parentID, "child").First(&child).Error; err != nil {
		http.Error(w, "Child not found or unauthorized", http.StatusNotFound)
		return
	}

	var readingLogs []models.ReadingLog
	if err := h.db.Where("child_id = ?", childID).Order("date DESC, created_at DESC").Find(&readingLogs).Error; err != nil {
		http.Error(w, "Failed to fetch reading logs", http.StatusInternalServerError)
		return
	}

	var responses []ReadingLogResponse
	for _, log := range readingLogs {
		response := ReadingLogResponse{
			ID:             log.ID,
			Title:          log.Title,
			Author:         log.Author,
			Status:         log.Status,
			Date:           log.Date,
			OpenLibraryKey: log.OpenLibraryKey,
			CoverID:        log.CoverID,
			CreatedAt:      log.CreatedAt,
		}
		responses = append(responses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
} 