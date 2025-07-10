package handlers

import (
	"github.com/lcleigh/page-hoppers-backend/models"
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "net/http"
)

// DB will be used to interact with the database
var DB *gorm.DB

// Parent Login Handler
func ParentLogin(c *gin.Context) {
    var loginDetails struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    // Parse the incoming JSON
    if err := c.ShouldBindJSON(&loginDetails); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var parent models.User
    // Look for the parent in the database with the provided credentials
    if err := DB.Where("username = ? AND password = ?", loginDetails.Username, loginDetails.Password).First(&parent).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Respond with the parent user info (could be a token in a real application)
    c.JSON(http.StatusOK, parent)
}

// Child Login Handler
func ChildLogin(c *gin.Context) {
    var loginDetails struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    // Parse the incoming JSON
    if err := c.ShouldBindJSON(&loginDetails); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var child models.User
    // Look for the child in the database with the provided credentials
    if err := DB.Where("username = ? AND password = ?", loginDetails.Username, loginDetails.Password).First(&child).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Respond with the child user info
    c.JSON(http.StatusOK, child)
}