package handlers

import (
	"github.com/lcleigh/page-hoppers-backend/models"
	"github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "net/http"
    "golang.org/x/crypto/bcrypt"
)

// DB will be used to interact with the database
var DB *gorm.DB

// Parent Login Handler
func ParentLogin(c *gin.Context) {
    var loginDetails struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    // Parse the incoming JSON
    if err := c.ShouldBindJSON(&loginDetails); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var parent models.User
    // Look for the parent in the database with the provided credentials
    if err := DB.Where("email = ? AND password = ?", loginDetails.Email, loginDetails.Password).First(&parent).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Respond with the parent user info (could be a token in a real application)
    c.JSON(http.StatusOK, parent)
}

// Child Login Handler
func ChildLogin(c *gin.Context) {
    var loginDetails struct {
        ChildID uint   `json:"childId"`
        PIN     string `json:"pin"`
    }

    if err := c.ShouldBindJSON(&loginDetails); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var child models.User
    if err := DB.Where("id = ? AND role = ?", loginDetails.ChildID, "child").First(&child).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Compare the hashed PIN
    if err := bcrypt.CompareHashAndPassword([]byte(child.PIN), []byte(loginDetails.PIN)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
        return
    }

    // Respond with the child user info (or a token, as needed)
    c.JSON(http.StatusOK, child)
}