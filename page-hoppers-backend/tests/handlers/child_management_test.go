package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"page-hoppers-backend/internal/models"
	"page-hoppers-backend/internal/handlers"
)

// TestCreateChildSuccess tests creating a child for a parent
func TestCreateChildSuccess(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a parent user
	parent := createTestParent(db, "Test Parent", "parent@example.com", "testpassword123")

	// Prepare child creation payload
	payload := handlers.CreateChildRequest{
		Name:     "Child One",
		Age:      8,
		PIN:      "1234",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	// Create request and set parent context
	req := httptest.NewRequest("POST", "/api/children", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "user_id", parent.ID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Call handler
	handler.CreateChild(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check response body
	var response handlers.ChildResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Username != payload.Username {
		t.Errorf("expected username %s, got %s", payload.Username, response.Username)
	}
	if response.ID == 0 {
		t.Error("expected non-zero child ID")
	}

	// Verify child was created in database
	var child models.User
	if err := db.Where("id = ?", response.ID).First(&child).Error; err != nil {
		t.Errorf("child was not created in database: %v", err)
	}
	if child.ParentID == nil || *child.ParentID != parent.ID {
		t.Errorf("expected parent ID %d, got %v", parent.ID, child.ParentID)
	}
	if child.Role != "child" {
		t.Errorf("expected role 'child', got %s", child.Role)
	}
}

// TestCreateChildValidation tests validation errors when creating a child
func TestCreateChildValidation(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := handlers.NewAuthHandler(db, []byte("test-secret"))

	// Create a parent user
	parent := createTestParent(db, "Test Parent", "parent@example.com", "testpassword123")

	// Missing username
	payload := handlers.CreateChildRequest{
		Name: "Child One", Age: 8, PIN: "1234",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/children", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "user_id", parent.ID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.CreateChild(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	// Missing PIN
	payload2 := handlers.CreateChildRequest{
		Username: "child2", Name: "Child Two", Age: 8,
	}
	body2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest("POST", "/api/children", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	ctx2 := context.WithValue(req2.Context(), "user_id", parent.ID)
	req2 = req2.WithContext(ctx2)
	w2 := httptest.NewRecorder()
	handler.CreateChild(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w2.Code)
	}

	// Age <= 0
	payload3 := handlers.CreateChildRequest{
		Username: "child3", Name: "Child Three", Age: 0, PIN: "1234",
	}
	body3, _ := json.Marshal(payload3)
	req3 := httptest.NewRequest("POST", "/api/children", bytes.NewReader(body3))
	req3.Header.Set("Content-Type", "application/json")
	ctx3 := context.WithValue(req3.Context(), "user_id", parent.ID)
	req3 = req3.WithContext(ctx3)
	w3 := httptest.NewRecorder()
	handler.CreateChild(w3, req3)
	if w3.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w3.Code)
	}
} 