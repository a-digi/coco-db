// Verwaltung von Datenbanken: Anlegen, Löschen, Auflisten
// Entspricht Punkt 3 der Roadmap und den Anforderungen aus PROJECT_REQUIREMENT.md

package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/a-digi/coco-db/src/response"
)

var dbNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

// DatabaseRequest repräsentiert das Request-Objekt für das Anlegen einer Datenbank
// (entspricht API-Design und PROJECT_REQUIREMENT.md)
type DatabaseRequest struct {
	Name string `json:"name"`
}

// DatabaseHandler verarbeitet alle Datenbank-Operationen (POST, DELETE, GET)
func DatabaseHandler(w http.ResponseWriter, r *http.Request) *response.APIResponse {
	dataDir := "./data" // TODO: Aus Konfiguration laden
	switch r.Method {
	case http.MethodPost:
		return handleCreateDatabase(w, r, dataDir)
	case http.MethodDelete:
		return handleDeleteDatabase(w, r, dataDir)
	case http.MethodGet:
		return handleListDatabases(w, r, dataDir)
	default:
		return response.WriteError(w, http.StatusMethodNotAllowed, "ERR_METHOD_NOT_ALLOWED", "Methode nicht erlaubt", "")
	}
}

// handleCreateDatabase legt eine neue Datenbank an
func handleCreateDatabase(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	var req DatabaseRequest
	if err := decodeJSON(r, &req); err != nil {
		return response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_JSON", err.Error(), "")
	}
	if !dbNamePattern.MatchString(req.Name) {
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

// handleDeleteDatabase löscht eine Datenbank
func handleDeleteDatabase(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	name := strings.TrimPrefix(r.URL.Path, "/api/databases/")
	if !dbNamePattern.MatchString(name) {
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

// handleListDatabases listet alle Datenbanken auf
func handleListDatabases(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	entries, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return response.WriteError(w, http.StatusInternalServerError, "ERR_IO", err.Error(), "")
	}
	dbs := []string{}
	for _, entry := range entries {
		if entry.IsDir() && dbNamePattern.MatchString(entry.Name()) {
			dbs = append(dbs, entry.Name())
		}
	}
	return response.WriteSuccess(w, map[string]interface{}{ "databases": dbs }, "")
}

// decodeJSON hilft beim Parsen von JSON-Requests
func decodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
