package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"page-hoppers-backend/internal/models"
)

// Add this method to your ReadingLogHandler
func (h *ReadingLogHandler) GetReadingSummary(c *gin.Context) {
	// Get the child ID from the URL
	childIDStr := c.Param("id")
	childIDUint, err := strconv.ParseUint(childIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid child ID"})
		return
	}
	childID := uint(childIDUint)

	// Make sure the child exists
	var child models.User
	if err := h.DB.Where("id = ? AND role = ?", childID, "child").First(&child).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Child not found"})
		return
	}

	// Fetch reading logs
	var logs []models.ReadingLog
	if err := h.DB.Where("child_id = ?", childID).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reading logs"})
		return
	}

	// Compute summary
	started := 0
	completed := 0
	for _, log := range logs {
		if log.Status == "started" {
			started++
		} else if log.Status == "completed" {
			completed++
		}
	}

	// Return JSON summary
	c.JSON(http.StatusOK, gin.H{
		"child_id":  childID,
		"name":      child.Name,
		"currentBook": gin.H{"title": "TEST BOOK"},
		"started":   started,
		"completed": completed,
		"total":     len(logs),
	})
}
