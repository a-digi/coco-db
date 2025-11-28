// Unit- und Integrationstests für das Servermodul
// Testet grundlegende Initialisierung und Health-Check-Endpoint

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler gab falschen Statuscode zurück: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler gab unerwarteten Body zurück: got %v want %v", rr.Body.String(), expected)
	}
}

