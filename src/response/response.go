package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse ist die Standardstruktur für alle API-Antworten.
type APIResponse struct {
	Success       bool        `json:"success"`
	Data          interface{} `json:"data,omitempty"`
	Error         *APIError   `json:"error,omitempty"`
	ExecutionTime string      `json:"executionTime,omitempty"`
}

// APIError beschreibt einen standardisierten Fehler.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteSuccess schreibt eine erfolgreiche Antwort mit Daten und gibt das APIResponse-Objekt zurück.
func WriteSuccess(w http.ResponseWriter, data interface{}, execTime string) *APIResponse {
	resp := &APIResponse{
		Success:       true,
		Data:          data,
		ExecutionTime: execTime,
	}
	resp.writeJSON(w, http.StatusOK)
	return resp
}

// WriteError schreibt eine Fehlerantwort mit Code und Nachricht und gibt das APIResponse-Objekt zurück.
func WriteError(w http.ResponseWriter, status int, code, message, execTime string) *APIResponse {
	resp := &APIResponse{
		Success:       false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
		ExecutionTime: execTime,
	}
	resp.writeJSON(w, status)
	return resp
}

// writeJSON serialisiert die Antwort als JSON und setzt die Header.
func (resp *APIResponse) writeJSON(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
