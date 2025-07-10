package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lcleigh/page-hoppers-backend/handlers"
)

// TestChildLoginSuccess tests successful child login
func TestChildLoginSuccess(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a parent and a child
	parent := createTestParent(db, "Test Parent", "parent@example.com", "parentpass")
	child := createTestChild(db, "childuser", "Child Name", 8, parent.ID, "1234")

	// Prepare login payload
	payload := handlers.ChildLoginRequest{
		ChildID: child.ID,
		PIN:    "1234",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.ChildLogin(w, req)

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

// TestChildLoginInvalidPIN tests login with wrong PIN
func TestChildLoginInvalidPIN(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a parent and a child
	parent := createTestParent(db, "Test Parent", "parent@example.com", "parentpass")
	child := createTestChild(db, "childuser", "Child Name", 8, parent.ID, "1234")

	// Prepare login payload with wrong PIN
	payload := handlers.ChildLoginRequest{
		ChildID: child.ID,
		PIN:    "9999", // Wrong PIN
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.ChildLogin(w, req)

	// Should get unauthorized status
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// TestChildLoginNonExistentChild tests login with child ID that doesn't exist
func TestChildLoginNonExistentChild(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Prepare login payload with non-existent child ID
	payload := handlers.ChildLoginRequest{
		ChildID: 999, // Non-existent child ID
		PIN:    "1234",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Create request
	req := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.ChildLogin(w, req)

	// Should get unauthorized status
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// TestChildLoginMissingFields tests validation of required fields
func TestChildLoginMissingFields(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Test with missing ChildID
	payload1 := map[string]string{
		"pin": "1234",
		// Missing childId
	}
	body1, _ := json.Marshal(payload1)
	req1 := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()

	handler.ChildLogin(w1, req1)

	// Should get unauthorized status (empty childId treated as invalid credentials)
	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w1.Code)
	}

	// Test with missing PIN
	payload2 := map[string]interface{}{
		"childId": 1,
		// Missing pin
	}
	body2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	handler.ChildLogin(w2, req2)

	// Should get unauthorized status (empty pin treated as invalid credentials)
	if w2.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w2.Code)
	}
}

// TestChildLoginParentUser tests that parent users cannot login through child endpoint
func TestChildLoginParentUser(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a parent user
	parent := createTestParent(db, "Test Parent", "parent@example.com", "parentpass")

	// Try to login as child with parent ID
	payload := handlers.ChildLoginRequest{
		ChildID: parent.ID,
		PIN:    "parentpass",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ChildLogin(w, req)

	// Should get unauthorized status because parent role doesn't match child
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// TestChildLoginEmptyPIN tests login with empty PIN
func TestChildLoginEmptyPIN(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a parent and a child
	parent := createTestParent(db, "Test Parent", "parent@example.com", "parentpass")
	child := createTestChild(db, "childuser", "Child Name", 8, parent.ID, "1234")

	// Prepare login payload with empty PIN
	payload := handlers.ChildLoginRequest{
		ChildID: child.ID,
		PIN:    "", // Empty PIN
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ChildLogin(w, req)

	// Should get unauthorized status
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// TestChildLoginMalformedJSON tests handling of malformed JSON
func TestChildLoginMalformedJSON(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create request with malformed JSON
	req := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader([]byte(`{"childId": 1, "pin": "1234"`))) // Missing closing brace
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ChildLogin(w, req)

	// Should get bad request status
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestChildLoginMultipleChildren tests login with different children
func TestChildLoginMultipleChildren(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a parent and multiple children
	parent := createTestParent(db, "Test Parent", "parent@example.com", "parentpass")
	child1 := createTestChild(db, "child1", "Child One", 8, parent.ID, "1111")
	child2 := createTestChild(db, "child2", "Child Two", 10, parent.ID, "2222")

	// Test login for first child
	payload1 := handlers.ChildLoginRequest{
		ChildID: child1.ID,
		PIN:    "1111",
	}
	body1, _ := json.Marshal(payload1)
	req1 := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()

	handler.ChildLogin(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("expected status 200 for child1, got %d", w1.Code)
	}

	// Test login for second child
	payload2 := handlers.ChildLoginRequest{
		ChildID: child2.ID,
		PIN:    "2222",
	}
	body2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest("POST", "/api/auth/child/login", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	handler.ChildLogin(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected status 200 for child2, got %d", w2.Code)
	}

	// Verify both tokens are different
	var response1, response2 handlers.LoginResponse
	json.Unmarshal(w1.Body.Bytes(), &response1)
	json.Unmarshal(w2.Body.Bytes(), &response2)

	if response1.Token == response2.Token {
		t.Error("expected different tokens for different children")
	}
} 