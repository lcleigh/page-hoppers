package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lcleigh/page-hoppers-backend/models"
	"github.com/lcleigh/page-hoppers-backend/handlers"
)

// TestParentRegister tests the parent registration endpoint
func TestParentRegister(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Test data
	payload := handlers.ParentRegisterRequest{
		Name:     "Test Parent",
		Email:    "parent@example.com",
		Password: "testpassword123",
	}
	
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/api/auth/parent/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.ParentRegister(w, req)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["message"] != "Parent registered successfully" {
		t.Errorf("expected success message, got %s", response["message"])
	}

	// Verify user was created in database
	var user models.User
	if err := db.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		t.Errorf("user was not created in database: %v", err)
	}

	if user.Role != "parent" {
		t.Errorf("expected role 'parent', got %s", user.Role)
	}

	if user.Username != payload.Name {
		t.Errorf("expected username %s, got %s", payload.Name, user.Username)
	}
}

// TestParentRegisterDuplicateEmail tests that duplicate emails are rejected
func TestParentRegisterDuplicateEmail(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create first user
	payload1 := handlers.ParentRegisterRequest{
		Name:     "Test Parent 1",
		Email:    "parent@example.com",
		Password: "testpassword123",
	}
	
	body1, _ := json.Marshal(payload1)
	req1 := httptest.NewRequest("POST", "/api/auth/parent/register", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.ParentRegister(w1, req1)

	// Try to create second user with same email
	payload2 := handlers.ParentRegisterRequest{
		Name:     "Test Parent 2",
		Email:    "parent@example.com", // Same email
		Password: "testpassword456",
	}
	
	body2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest("POST", "/api/auth/parent/register", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.ParentRegister(w2, req2)

	// Should get conflict status
	if w2.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w2.Code)
	}
}

// TestParentRegisterMissingFields tests validation of required fields
func TestParentRegisterMissingFields(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Test with missing email
	payload := map[string]string{
		"name":     "Test Parent",
		"password": "testpassword123",
		// Missing email
	}
	
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/parent/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ParentRegister(w, req)

	// Should get bad request status
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestParentLogin tests successful parent login
func TestParentLogin(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// First, create a parent user
	createTestParent(db, "Test Parent", "parent@example.com", "testpassword123")

	// Test login with correct credentials
	payload := handlers.ParentLoginRequest{
		Email:    "parent@example.com",
		Password: "testpassword123",
	}
	
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/api/auth/parent/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.ParentLogin(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check response body contains token
	var response handlers.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Token == "" {
		t.Error("expected token in response, got empty string")
	}

	// Verify token is valid JWT
	if len(response.Token) < 50 { // Basic check that token looks like JWT
		t.Errorf("token seems too short: %s", response.Token)
	}
}

// TestParentLoginInvalidCredentials tests login with wrong password
func TestParentLoginInvalidCredentials(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a parent user
	createTestParent(db, "Test Parent", "parent@example.com", "testpassword123")

	// Test login with wrong password
	payload := handlers.ParentLoginRequest{
		Email:    "parent@example.com",
		Password: "wrongpassword",
	}
	
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/parent/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ParentLogin(w, req)

	// Should get unauthorized status
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// TestParentLoginNonExistentEmail tests login with email that doesn't exist
func TestParentLoginNonExistentEmail(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Test login with non-existent email
	payload := handlers.ParentLoginRequest{
		Email:    "nonexistent@example.com",
		Password: "testpassword123",
	}
	
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/parent/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ParentLogin(w, req)

	// Should get unauthorized status
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// TestParentLoginMissingFields tests validation of required fields
func TestParentLoginMissingFields(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Test with missing email
	payload := map[string]string{
		"password": "testpassword123",
		// Missing email
	}
	
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/parent/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ParentLogin(w, req)

	// Should get unauthorized status (empty email treated as invalid credentials)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	// Test with missing password
	payload2 := map[string]string{
		"email": "parent@example.com",
		// Missing password
	}
	
	body2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest("POST", "/api/auth/parent/login", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	handler.ParentLogin(w2, req2)

	// Should get unauthorized status (empty password treated as invalid credentials)
	if w2.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w2.Code)
	}
}

// TestParentLoginChildUser tests that child users cannot login through parent endpoint
func TestParentLoginChildUser(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a child user with email (unusual but possible)
	child := &models.User{
		Username: "Test Child",
		Email:    "child@example.com",
		Password: hashPassword("testpassword123"),
		Role:     "child",
	}
	db.Create(child)

	// Try to login as parent with child credentials
	payload := handlers.ParentLoginRequest{
		Email:    "child@example.com",
		Password: "testpassword123",
	}
	
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/parent/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ParentLogin(w, req)

	// Should get unauthorized status because child role doesn't match parent
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
} 