// Verwaltung von Datenbanken: Anlegen, Löschen, Auflisten
// Entspricht Punkt 3 der Roadmap und den Anforderungen aus PROJECT_REQUIREMENT.md

package database

import (
	"net/http"
	"regexp"
	"github.com/a-digi/coco-db/src/response"
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
