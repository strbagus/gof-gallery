package event_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/strbagus/gof-gallery/internal/database"
	"github.com/strbagus/gof-gallery/module/event"
	res "github.com/strbagus/gof-gallery/pkg/response"
)

func TestEventTokensEndpoints(t *testing.T) {
	// Load environment variables
	_ = godotenv.Load("../../.env")

	// Set APP_ENV to development to bypass AdminMiddleware
	os.Setenv("APP_ENV", "development")

	// Initialize Postgres
	database.InitPostgres()
	defer database.ClosePostgres()

	ctx := context.Background()

	// Clean up any pre-existing test data
	_, _ = database.PgxPool.Exec(ctx, "DELETE FROM public.events WHERE slug = $1", "test-event-token-slug")

	// 1. Create a test event
	createReq := &event.CreateRequest{
		Event: event.Event{
			Name:      "Test Event for Tokens",
			Slug:      "test-event-token-slug",
			IsPrivate: true,
		},
	}
	err := event.CreateEvent(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}
	defer func() {
		_, _ = database.PgxPool.Exec(ctx, "DELETE FROM public.events WHERE slug = $1", "test-event-token-slug")
	}()

	// Retrieve event detail to get ID
	evt, err := event.GetEventBySlug(ctx, "test-event-token-slug")
	if err != nil {
		t.Fatalf("Failed to retrieve test event: %v", err)
	}
	_ = evt

	// 2. Setup Fiber application
	app := fiber.New()
	event.RegisterRoutes(app)

	// --- Endpoint 1: POST /events/:slug/tokens ---
	payload := map[string]any{
		"duration_days": 10,
	}
	payloadBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/events/test-event-token-slug/tokens", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var postResponse res.JSONResponse
	if err := json.NewDecoder(resp.Body).Decode(&postResponse); err != nil {
		t.Fatalf("Failed to decode POST response: %v", err)
	}

	if postResponse.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", postResponse.Status)
	}

	// Extract generated token
	dataMap, ok := postResponse.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected Data to be map[string]any, got %T", postResponse.Data)
	}

	tokenStr, _ := dataMap["token"].(string)
	if len(tokenStr) != 64 {
		t.Errorf("Expected token length 64, got %d", len(tokenStr))
	}

	// --- Endpoint 2: GET /events/:slug/tokens ---
	reqGet := httptest.NewRequest("GET", "/events/test-event-token-slug/tokens", nil)
	respGet, err := app.Test(reqGet)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if respGet.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", respGet.StatusCode)
	}

	var getResponse res.JSONResponse
	if err := json.NewDecoder(respGet.Body).Decode(&getResponse); err != nil {
		t.Fatalf("Failed to decode GET response: %v", err)
	}

	tokensList, ok := getResponse.Data.([]any)
	if !ok || len(tokensList) == 0 {
		t.Fatalf("Expected non-empty list of tokens, got: %v", getResponse.Data)
	}

	// --- Endpoint 3: DELETE /events/:slug/tokens/:token ---
	reqDelete := httptest.NewRequest("DELETE", "/events/test-event-token-slug/tokens/"+tokenStr, nil)
	respDelete, err := app.Test(reqDelete)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}

	if respDelete.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", respDelete.StatusCode)
	}

	// Verify token is deleted
	reqGetAfter := httptest.NewRequest("GET", "/events/test-event-token-slug/tokens", nil)
	respGetAfter, err := app.Test(reqGetAfter)
	if err != nil {
		t.Fatalf("GET after DELETE request failed: %v", err)
	}

	var getResponseAfter res.JSONResponse
	_ = json.NewDecoder(respGetAfter.Body).Decode(&getResponseAfter)
	tokensListAfter, _ := getResponseAfter.Data.([]any)
	if len(tokensListAfter) != 0 {
		t.Errorf("Expected token list to be empty after delete, got %d items", len(tokensListAfter))
	}

	// --- Endpoint 4: PUT /events/:slug ---
	updatePayload := map[string]any{
		"name":       "Updated Event Name",
		"slug":       "test-event-token-slug-updated",
		"location":   "Updated Location",
		"is_private": false,
	}
	updateBytes, _ := json.Marshal(updatePayload)
	reqUpdate := httptest.NewRequest("PUT", "/events/test-event-token-slug", bytes.NewReader(updateBytes))
	reqUpdate.Header.Set("Content-Type", "application/json")
	respUpdate, err := app.Test(reqUpdate)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	if respUpdate.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", respUpdate.StatusCode)
	}

	// Verify updated event details
	updatedEvt, err := event.GetEventBySlug(ctx, "test-event-token-slug-updated")
	if err != nil {
		t.Fatalf("Failed to retrieve updated event: %v", err)
	}
	if updatedEvt.Name != "Updated Event Name" {
		t.Errorf("Expected name 'Updated Event Name', got '%s'", updatedEvt.Name)
	}

	// Clean up using the updated slug later
	defer func() {
		_, _ = database.PgxPool.Exec(ctx, "DELETE FROM public.events WHERE slug = $1", "test-event-token-slug-updated")
	}()

	// --- Endpoint 5: DELETE /events/:slug ---
	reqDeleteEvt := httptest.NewRequest("DELETE", "/events/test-event-token-slug-updated", nil)
	respDeleteEvt, err := app.Test(reqDeleteEvt)
	if err != nil {
		t.Fatalf("DELETE event request failed: %v", err)
	}
	if respDeleteEvt.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", respDeleteEvt.StatusCode)
	}

	// Verify it is soft deleted (deleted_at is set, GetEventBySlug returns ErrNoRows)
	_, err = event.GetEventBySlug(ctx, "test-event-token-slug-updated")
	if err == nil {
		t.Errorf("Expected event to be soft-deleted and not found, but it was found")
	}
}
