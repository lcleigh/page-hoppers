package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
	"github.com/lcleigh/page-hoppers-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db     *gorm.DB
	secret []byte
}

func NewAuthHandler(db *gorm.DB, secret []byte) *AuthHandler {
	return &AuthHandler{
		db:     db,
		secret: secret,
	}
}

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

func (h *AuthHandler) ParentLogin(w http.ResponseWriter, r *http.Request) {
	var req ParentLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.Where("email = ? AND role = ?", req.Email, "parent").First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    "parent",
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.secret)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}

func (h *AuthHandler) ChildLogin(w http.ResponseWriter, r *http.Request) {
	var req ChildLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var child models.User
	if err := h.db.Where("id = ? AND role = ?", req.ChildID, "child").First(&child).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(child.PIN), []byte(req.PIN)); err != nil {
		http.Error(w, "Invalid PIN", http.StatusUnauthorized)
		return
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   child.ID,
		"parent_id": child.ParentID,
		"role":      "child",
		"exp":       time.Now().Add(time.Hour * 12).Unix(), // Shorter expiration for child tokens
	})

	tokenString, err := token.SignedString(h.secret)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}

func (h *AuthHandler) GetChildren(w http.ResponseWriter, r *http.Request) {
	// Extract parent ID from JWT token
	parentID := r.Context().Value("user_id").(uint)

	var children []models.User
	if err := h.db.Where("parent_id = ? AND role = ?", parentID, "child").Find(&children).Error; err != nil {
		http.Error(w, "Failed to fetch children", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(children)
}

type CreateChildRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	PIN      string `json:"pin"`
}

type ChildResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func (h *AuthHandler) CreateChild(w http.ResponseWriter, r *http.Request) {
	parentID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateChildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.PIN == "" || req.Name == "" || req.Age <= 0 {
		http.Error(w, "Username, Name, Age, and PIN are required", http.StatusBadRequest)
		return
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash PIN", http.StatusInternalServerError)
		return
	}

	child := models.User{
		Username: req.Username,
		Name:     req.Name,
		Age:      req.Age,
		PIN:      string(hashedPIN),
		Role:     "child",
		ParentID: &parentID,
	}

	if err := h.db.Create(&child).Error; err != nil {
		http.Error(w, "Failed to create child", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ChildResponse{ID: child.ID, Username: child.Username})
}

type ParentRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) ParentRegister(w http.ResponseWriter, r *http.Request) {
	var req ParentRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Name, email, and password are required", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create parent user
	parent := models.User{
		Username: req.Name, // Use name as username for now
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "parent",
	}

	if err := h.db.Create(&parent).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Parent registered successfully"})
} 