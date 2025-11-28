package database

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-digi/coco-db/src/response"
)

func HandleDeleteDatabase(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	return handleDeleteDatabase(w, r, dataDir)
}

func handleDeleteDatabase(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	// REST: /api/databases/{dbname}
	prefix := "/api/databases/"
	if !strings.HasPrefix(r.URL.Path, prefix) || len(r.URL.Path) <= len(prefix) {
		return response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
	}
	name := r.URL.Path[len(prefix):]
	if len(name) == 0 || !DbNamePattern.MatchString(name) {
		return response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
	}
	dbPath := filepath.Join(dataDir, name)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return response.WriteError(w, http.StatusNotFound, "ERR_DB_NOT_FOUND", "Datenbank nicht gefunden", "")
	}
	if err := os.RemoveAll(dbPath); err != nil {
		return response.WriteError(w, http.StatusInternalServerError, "ERR_IO", err.Error(), "")
	}
	return response.WriteSuccess(w, map[string]string{"name": name}, "")
}
