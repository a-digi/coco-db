package database

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-digi/coco-db/src/response"
)

type UpdateDatabaseRequest struct {
	NewName string `json:"newName"`
}

// Platzhalter für zukünftige Update-Logik
func handleUpdateDatabase(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	oldName := strings.TrimPrefix(r.URL.Path, "/api/databases/update/")
	if len(oldName) == 0 || !DbNamePattern.MatchString(oldName) {
		return response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger alter Datenbankname", "")
	}
	var req UpdateDatabaseRequest
	if err := decodeJSON(r, &req); err != nil {
		return response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_JSON", err.Error(), "")
	}
	if !DbNamePattern.MatchString(req.NewName) {
		return response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger neuer Datenbankname", "")
	}
	if req.NewName == oldName {
		return response.WriteError(w, http.StatusConflict, "ERR_DB_SAME_NAME", "Neuer Name ist identisch mit altem Namen", "")
	}
	oldPath := filepath.Join(dataDir, oldName)
	newPath := filepath.Join(dataDir, req.NewName)
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return response.WriteError(w, http.StatusNotFound, "ERR_DB_NOT_FOUND", "Alte Datenbank nicht gefunden", "")
	}
	if _, err := os.Stat(newPath); err == nil {
		return response.WriteError(w, http.StatusConflict, "ERR_DB_EXISTS", "Ziel-Datenbankname existiert bereits", "")
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, "ERR_IO", err.Error(), "")
	}
	return response.WriteSuccess(w, map[string]string{"oldName": oldName, "newName": req.NewName}, "")
}

func HandleUpdateDatabase(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	return handleUpdateDatabase(w, r, dataDir)
}
