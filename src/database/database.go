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

// Chain of Responsibility für GET-Routen

type DBHandlerFunc func(w http.ResponseWriter, r *http.Request, dataDir string) (handled bool, resp *response.APIResponse)

type DBHandlerChain struct {
	handlers []DBHandlerFunc
	dataDir  string
}

func NewDBHandlerChain(dataDir string) *DBHandlerChain {
	return &DBHandlerChain{
		handlers: []DBHandlerFunc{
			routeHandler("/api/databases", handleListDatabases),
			routeHandler("/api/databases/create/", handleCreateDatabase),
			routeHandler("/api/databases/delete/", handleDeleteDatabase),
			routeHandler("/api/databases/update/", handleUpdateDatabase),
			routeHandlerRESTUpdate(), // <--- REST-Update-Handler für PUT /api/databases/{dbname}
		},
		dataDir: dataDir,
	}
}

// Handler für exakten oder Prefix-Match (nur GET)
func routeHandler(routePattern string, fn func(http.ResponseWriter, *http.Request, string) *response.APIResponse) DBHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, dataDir string) (bool, *response.APIResponse) {
		if routePattern == r.URL.Path || (routePattern[len(routePattern)-1] == '/' && len(r.URL.Path) > len(routePattern) && r.URL.Path[:len(routePattern)] == routePattern) {
			return true, fn(w, r, dataDir)
		}
		return false, nil
	}
}

// REST-Handler für PUT /api/databases/{dbname}
func routeHandlerRESTUpdate() DBHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, dataDir string) (bool, *response.APIResponse) {
		if r.Method != http.MethodPut {
			return false, nil
		}
		prefix := "/api/databases/"
		if !startsWith(r.URL.Path, prefix) {
			return false, nil
		}
		dbname := r.URL.Path[len(prefix):]
		if dbname == "" || !DbNamePattern.MatchString(dbname) {
			return false, nil
		}
		return true, handleUpdateDatabaseREST(w, r, dataDir, dbname)
	}
}

func (c *DBHandlerChain) Serve(w http.ResponseWriter, r *http.Request) *response.APIResponse {
	for _, h := range c.handlers {
		handled, resp := h(w, r, c.dataDir)
		if handled {
			return resp
		}
	}
	return response.WriteError(w, http.StatusNotFound, "ERR_ROUTE_NOT_FOUND", "Route nicht gefunden", "")
}

// Handler verarbeitet alle Datenbank-Operationen (nur GET) über die Chain
func Handler(w http.ResponseWriter, r *http.Request) *response.APIResponse {
	chain := NewDBHandlerChain("./data") // TODO: Aus Konfiguration laden
	return chain.Serve(w, r)
}

// REST-Logik für PUT /api/databases/{dbname}
func handleUpdateDatabaseREST(w http.ResponseWriter, r *http.Request, dataDir, oldName string) *response.APIResponse {
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

// Hilfsfunktion für Prefix-Matching
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
