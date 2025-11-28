package test

import (
	"net/http"
	"net/http/httptest"
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
	req := httptest.NewRequest(http.MethodGet, "/api/databases/create/testdb", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusConflict {
		t.Errorf("/api/databases/create/ Route gibt unerwarteten Status zurück: %d", w.Code)
	}
}

func TestRouter_DatabaseDeleteRoute(t *testing.T) {
	router := server.SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/databases/delete/testdb", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
		t.Errorf("/api/databases/delete/ Route gibt unerwarteten Status zurück: %d", w.Code)
	}
}

func TestRouter_DatabaseUpdateRoute(t *testing.T) {
	router := server.SetupRouter()
	updateBody := httptest.NewRequest(http.MethodGet, "/api/databases/update/testdb", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, updateBody)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusConflict && w.Code != http.StatusOK {
		t.Errorf("/api/databases/update/ Route gibt unerwarteten Status zurück: %d", w.Code)
	}
}

