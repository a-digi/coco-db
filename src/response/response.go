package response

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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
	if execTime != "" && strings.HasSuffix(execTime, "µs") {
		if dur, err := strconv.ParseFloat(strings.TrimSuffix(execTime, "µs"), 64); err == nil {
			execTime = strconv.FormatInt(int64(dur/1000), 10) + "ms"
		}
	}

	resp := &APIResponse{
		Success:       true,
		Data:          data,
		ExecutionTime: execTime,
		HttpCode:      http.StatusOK,
	}

	return resp
}

// WriteError schreibt eine Fehlerantwort mit Code und Nachricht und gibt das APIResponse-Objekt zurück.
func WriteError(w http.ResponseWriter, status int, code, message, execTime string) *APIResponse {
	if execTime != "" && strings.HasSuffix(execTime, "µs") {
		if dur, err := strconv.ParseFloat(strings.TrimSuffix(execTime, "µs"), 64); err == nil {
			execTime = strconv.FormatInt(int64(dur/1000), 10) + "ms"
		}
	}
	resp := WriteErrorInternal(status, code, message, execTime)
	if w != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
	}
	return resp
}

func WriteErrorInternal(status int, code, message, execTime string) *APIResponse {
	if execTime == "" {
		execTime = "0s"
	}
	if execTime != "" && strings.HasSuffix(execTime, "µs") {
		if dur, err := strconv.ParseFloat(strings.TrimSuffix(execTime, "µs"), 64); err == nil {
			execTime = strconv.FormatInt(int64(dur/1000), 10) + "ms"
		}
	}
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
