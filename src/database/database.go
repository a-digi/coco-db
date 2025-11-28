// Verwaltung von Datenbanken: Anlegen, Löschen, Auflisten
// Entspricht Punkt 3 der Roadmap und den Anforderungen aus PROJECT_REQUIREMENT.md

package database

import (
	"net/http"
	"regexp"
	"github.com/a-digi/coco-db/src/response"
	"os"
	"path/filepath"
)

var DbNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

// DatabaseRequest repräsentiert das Request-Objekt für das Anlegen einer Datenbank
type DatabaseRequest struct {
	Name string `json:"name"`
}

// REST-Logik für PUT /api/databases/{dbname}
func handleUpdateDatabaseREST(w http.ResponseWriter, r *http.Request, dataDir, oldName string) *response.APIResponse {

	var req UpdateDatabaseRequest
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

// Hilfsfunktion für Prefix-Matching
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
