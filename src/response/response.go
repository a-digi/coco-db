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
	HttpCode      int         `json:"httpCode"`
}

// APIError beschreibt einen standardisierten Fehler.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteSuccess schreibt eine erfolgreiche Antwort mit Daten und gibt das APIResponse-Objekt zurück.
func WriteSuccess(data interface{}, execTime string) *APIResponse {
	resp := &APIResponse{
		Success:       true,
		Data:          data,
		ExecutionTime: execTime,
		HttpCode:      http.StatusOK,
	}

	return resp
}

// WriteError schreibt eine Fehlerantwort mit Code und Nachricht und gibt das APIResponse-Objekt zurück.
func WriteError(status int, code, message, execTime string) *APIResponse {
	resp := &APIResponse{
		Success:       false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
		ExecutionTime: execTime,
		HttpCode:      status,
	}

	return resp
}
