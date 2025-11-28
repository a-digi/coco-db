// Unit- und Integrationstests für das Servermodul
// Testet grundlegende Initialisierung und Health-Check-Endpoint

package server

import (
	"encoding/json"
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

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Antwort ist kein valides JSON: %v", err)
	}
	if !resp.Success || resp.Data.Status != "ok" {
		t.Errorf("Handler gab unerwarteten Body zurück: got %v", rr.Body.String())
	}
}
