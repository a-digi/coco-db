// Routing und HTTP-Handler für die Kernoperationen (CRUD, Index, Events, Transaktionen)
// Hier werden die Endpunkte gemäß API-Design angebunden

package server

import (
	"net/http"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/database"
	"github.com/a-digi/coco-db/src/table"
	"github.com/a-digi/coco-db/src/logger"
	paramrouter "github.com/a-digi/coco-db/src/server/router"
)

// apiHandler ist ein Wrapper, der Handler mit APIResponse-Signatur in http.HandlerFunc umwandelt
func apiHandler(fn func(http.ResponseWriter, *http.Request) *response.APIResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := fn(w, r)
		if resp == nil {
			response.WriteError(w, http.StatusInternalServerError, "internal_error", "Leere APIResponse", "")
		}
	}
}

// SetupRouter initialisiert das Routing für alle Kernoperationen
func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", apiHandler(HealthHandler))

	pr := paramrouter.NewParamRouter()

	// Datenbank-Endpunkte
	pr.HandleFunc("GET", "/api/databases", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		database.HandleListDatabases(w, r, "./data")
	})

	pr.HandleFunc("POST", "/api/databases", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		database.HandleCreateDatabase(w, r, "./data")
	})

	pr.HandleFunc("DELETE", "/api/databases/{dbname}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		database.HandleDeleteDatabase(w, r, "./data")
	})

	pr.HandleFunc("PUT", "/api/databases/{dbname}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		database.HandleUpdateDatabase(w, r, "./data")
	})

	// Tabellen-Endpunkt: POST /api/databases/{dbname}/tables
	pr.HandleFunc("POST", "/api/databases/{dbname}/tables", func(w http.ResponseWriter, r *http.Request, params map[string]string) {

		tc := &table.TableCreator{
			DataDir: "./data", // TODO: Aus config.json laden
			Logger:  &logger.NoopLogger{},
		}

		// Setze den Datenbanknamen als Header, damit die Handler-Logik konsistent bleibt
		if dbname, ok := params["dbname"]; ok {
			r.Header.Set("X-DB-Name", dbname)
		}
		tc.HandleCreateTable(w, r)
	})
	// Tabellen-Endpunkt: GET /api/databases/{dbname}/tables
	pr.HandleFunc("GET", "/api/databases/{dbname}/tables", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		listHandler := &table.ListTablesHandler{
			ResponseWriter: w,
			DataDir: "./data", // TODO: Aus config.json laden
			Logger:  &logger.NoopLogger{},
		}
		listHandler.HandleListTables(params["dbname"])
	})
	// Hier können weitere Table-Endpunkte ergänzt werden

	// Kombiniere Health-Mux und ParamRouter
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			mux.ServeHTTP(w, r)
			return
		}
		pr.ServeHTTP(w, r)
	})
}

// startsWith prüft, ob s mit prefix beginnt
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// indexOf gibt den Index des ersten Vorkommens von sep in s zurück, oder -1
func indexOf(s, sep string) int {

	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return i
		}
	}

	return -1
}
