package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"page-hoppers-backend/internal/models"
)

type AuthHandler struct {
	DB     *gorm.DB
	Secret []byte
}

func NewAuthHandler(db *gorm.DB, secret []byte) *AuthHandler {
	return &AuthHandler{
		DB:     db,
		Secret: secret,
	}
}

// Request/Response structs
type ParentLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChildLoginRequest struct {
	ChildID uint   `json:"childId"`
	PIN     string `json:"pin"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateChildRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	PIN  string `json:"pin"`
}

type ChildResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ParentRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ---------------------------
// Parent login
func (h *AuthHandler) ParentLogin(c *gin.Context) {
	var req ParentLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var parent models.User
	if err := h.DB.Where("email = ? AND role = ?", req.Email, "parent").First(&parent).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(parent.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": parent.ID,
		"role":    "parent",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(h.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

// ---------------------------
// Child login
func (h *AuthHandler) ChildLogin(c *gin.Context) {
	var req ChildLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var child models.User
	if err := h.DB.Where("id = ? AND role = ?", req.ChildID, "child").First(&child).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(child.PIN), []byte(req.PIN)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   child.ID,
		"parent_id": child.ParentID,
		"role":      "child",
		"exp":       time.Now().Add(12 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(h.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

// ---------------------------
// Get children for a parent
func (h *AuthHandler) GetChildren(c *gin.Context) {
	parentIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	parentID := parentIDValue.(uint)

	var children []models.User
	if err := h.DB.Where("parent_id = ? AND role = ?", parentID, "child").Find(&children).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch children"})
		return
	}

	c.JSON(http.StatusOK, children)
}

// ---------------------------
// Create a child
func (h *AuthHandler) CreateChild(c *gin.Context) {
	parentIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	parentID := parentIDValue.(uint)

	var req CreateChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Name == "" || req.Age <= 0 || req.PIN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, Age, and PIN are required"})
		return
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash PIN"})
		return
	}

	child := models.User{
		Name:     req.Name,
		Age:      req.Age,
		PIN:      string(hashedPIN),
		Role:     "child",
		ParentID: &parentID,
	}

	if err := h.DB.Create(&child).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create child"})
		return
	}

	c.JSON(http.StatusOK, ChildResponse{ID: child.ID, Name: child.Name})
}

// ---------------------------
// Parent registration
func (h *AuthHandler) ParentRegister(c *gin.Context) {
	var req ParentRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email, and password are required"})
		return
	}

	var existing models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	parent := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "parent",
	}

	if err := h.DB.Create(&parent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create parent"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Parent registered successfully"})
}