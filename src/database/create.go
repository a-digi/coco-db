package database

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-digi/coco-db/src/response"
)

func HandleCreateDatabase(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	return handleCreateDatabase(w, r, dataDir)
}

func handleCreateDatabase(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	var req DatabaseRequest
	if err := decodeJSON(r, &req); err != nil {
		return response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_JSON", err.Error(), "")
	}
	if !DbNamePattern.MatchString(req.Name) {
		return response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
	}
	dbPath := filepath.Join(dataDir, req.Name)
	if _, err := os.Stat(dbPath); err == nil {
		return response.WriteError(w, http.StatusConflict, "ERR_DB_EXISTS", "Datenbank existiert bereits", "")
	}
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, "ERR_IO", err.Error(), "")
	}
	return response.WriteSuccess(w, map[string]string{"name": req.Name}, "")
}

func decodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
