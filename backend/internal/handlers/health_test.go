package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealth(t *testing.T) {
	app := fiber.New()
	app.Get("/health", Health("1.0.0"))

	req := httptest.NewRequest("GET", "/health", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var body HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", body.Status)
	}

	if body.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", body.Version)
	}
}

func TestReady(t *testing.T) {
	app := fiber.New()
	app.Get("/ready", Ready())

	req := httptest.NewRequest("GET", "/ready", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var body ReadyResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body.Status != "ready" {
		t.Errorf("Expected status 'ready', got '%s'", body.Status)
	}

	if !body.Ready {
		t.Error("Expected ready to be true")
	}
}
