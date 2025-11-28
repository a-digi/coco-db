package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-digi/coco-db/src/server"
)

func TestRouter_HealthRoute(t *testing.T) {
	router := server.SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("/health Route gibt falschen Status zurück: %d", w.Code)
	}
}

func TestRouter_DatabasesRoute(t *testing.T) {
	router := server.SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/databases", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("/api/databases Route gibt unerwarteten Status zurück: %d", w.Code)
	}
}

func TestRouter_DatabaseCreateRoute(t *testing.T) {
	router := server.SetupRouter()
	// Nutze POST /api/databases mit JSON-Body
	body := strings.NewReader(`{"name":"testdb"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/databases", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusConflict {
		t.Errorf("/api/databases (POST) Route gibt unerwarteten Status zurück: %d", w.Code)
	}
}

func TestRouter_DatabaseDeleteRoute(t *testing.T) {
	router := server.SetupRouter()
	// Nutze DELETE /api/databases/testdb
	req := httptest.NewRequest(http.MethodDelete, "/api/databases/testdb", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
		t.Errorf("/api/databases/{dbname} (DELETE) Route gibt unerwarteten Status zurück: %d", w.Code)
	}
}

func TestRouter_DatabaseUpdateRoute(t *testing.T) {
	router := server.SetupRouter()
	// Nutze PUT /api/databases/testdb
	updateBody := strings.NewReader(`{"name":"testdb"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/databases/testdb", updateBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusConflict && w.Code != http.StatusOK {
		t.Errorf("/api/databases/{dbname} (PUT) Route gibt unerwarteten Status zurück: %d", w.Code)
	}
}
