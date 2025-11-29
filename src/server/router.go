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
	"encoding/json"
	"github.com/a-digi/coco-db/src/table/fields"
	entries "github.com/a-digi/coco-db/src/table/entries"
	"time"
	"github.com/a-digi/coco-db/src/query"
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
		dl := &database.DatabaseList{DataDir: "./data", Logger: &logger.NoopLogger{}}
		resp := dl.HandleListDatabases()
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})

	pr.HandleFunc("POST", "/api/databases", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		dbName := r.URL.Query().Get("name")
		if dbName == "" {
			var body struct{ Name string `json:"name"` }
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				dbName = body.Name
			}
		}
		dc := &database.DatabaseCreate{
			DataDir: "./data",
			Logger:  &logger.NoopLogger{},
		}
		resp := dc.HandleCreateDatabase(dbName)
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})

	pr.HandleFunc("DELETE", "/api/databases/{dbname}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		dbName := params["dbname"]
		dd := &database.DatabaseDelete{DataDir: "./data", Logger: &logger.NoopLogger{}}
		resp := dd.HandleDeleteDatabase(dbName)
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})

	pr.HandleFunc("PUT", "/api/databases/{dbname}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		oldName := params["dbname"]
		var req struct{ NewName string `json:"newName"` }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_DB_INVALID_JSON", "Ungültiges JSON: "+err.Error(), "")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		du := &database.DatabaseUpdate{SataDir: "./data"}
		resp := du.HandleUpdateDatabase(oldName, req.NewName)
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})

	// Tabellen-Endpunkt: POST /api/databases/{dbname}/tables
	pr.HandleFunc("POST", "/api/databases/{dbname}/tables", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		tc := &table.TableCreator{
			DataDir: "./data", // TODO: Aus config.json laden
			Logger:  &logger.NoopLogger{},
		}
		dbname := params["dbname"]
		var meta fields.TableMeta
		if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
			resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_JSON", "Ungültiges JSON: "+err.Error(), "")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		resp := tc.HandleCreateTable(dbname, meta)
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})
	// Tabellen-Endpunkt: GET /api/databases/{dbname}/tables
	pr.HandleFunc("GET", "/api/databases/{dbname}/tables", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		h := &table.ListTablesHandler{
			DataDir: "./data", // TODO: Aus config.json laden
			Logger:  &logger.NoopLogger{},
		}
		resp := h.HandleListTables(params["dbname"])
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})
	// Tabellen-Endpunkt: PUT /api/databases/{dbname}/tables/{tablename}
	pr.HandleFunc("PUT", "/api/databases/{dbname}/tables/{tablename}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		tu := &table.TableUpdate{
			DataDir: "./data",
			Logger:  &logger.NoopLogger{},
		}
		dbname := params["dbname"]
		tablename := params["tablename"]
		var meta fields.TableMeta
		if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
			resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_JSON", "Ungültiges JSON: "+err.Error(), "")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		resp := tu.HandleEditTable(dbname, tablename, meta)
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})
	// Tabellen-Endpunkt: DELETE /api/databases/{dbname}/tables/{tablename}
	pr.HandleFunc("DELETE", "/api/databases/{dbname}/tables/{tablename}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		td := &table.TableDelete{
			DataDir: "./data",
			Logger:  &logger.NoopLogger{},
		}
		dbname := params["dbname"]
		tablename := params["tablename"]
		resp := td.HandleDeleteTable(dbname, tablename)
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})
	// Tabellen-Endpunkt: GET /api/databases/{dbname}/tables/{tablename}
	pr.HandleFunc("GET", "/api/databases/{dbname}/tables/{tablename}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		tlm := &table.TableListMeta{
			DataDir: "./data",
			Logger:  &logger.NoopLogger{},
		}
		dbname := params["dbname"]
		tablename := params["tablename"]
		resp := tlm.HandleGetTable(dbname, tablename)
		w.Header().Set("Content-Type", "application/json")
		if resp != nil {
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
		}
	})
	// Einfüge-Endpunkt: POST /api/databases/{dbname}/tables/{tablename}/entries
	pr.HandleFunc("POST", "/api/databases/{dbname}/tables/{tablename}/entries", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		dbname := params["dbname"]
		tablename := params["tablename"]
		var entry map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_JSON", "Ungültiges JSON: "+err.Error(), "")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		resp := entries.InsertEntry(dbname, tablename, entry, "./data", &logger.NoopLogger{})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
	})

	// Edit-Endpunkt: PUT /api/databases/{dbname}/tables/{tablename}/entries/{entryid}
	pr.HandleFunc("PUT", "/api/databases/{dbname}/tables/{tablename}/entries/{entryid}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		dbname := params["dbname"]
		tablename := params["tablename"]
		entryid := params["entryid"]
		var entry map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_JSON", "Ungültiges JSON: "+err.Error(), "")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		editor := &entries.EntryEditor{DataDir: "./data", Logger: &logger.NoopLogger{}}
		err := editor.EditEntry(dbname, tablename, entryid, entry)
		var resp *response.APIResponse
		if err != nil {
			if err.Error() == "Eintrag nicht gefunden" {
				resp = response.WriteErrorInternal(404, "ERR_ENTRY_NOT_FOUND", err.Error(), "")
			} else if err.Error() == "Keine aktive Version vorhanden" {
				resp = response.WriteErrorInternal(404, "ERR_NO_ACTIVE_VERSION", err.Error(), "")
			} else {
				resp = response.WriteErrorInternal(400, "ERR_UPDATE_ENTRY", err.Error(), "")
			}
		} else {
			resp = response.WriteSuccess(map[string]interface{}{
				"entryId": entryid,
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			}, "")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
	})

	// Delete-Endpunkt: DELETE /api/databases/{dbname}/tables/{tablename}/entries/{entryid}
	pr.HandleFunc("DELETE", "/api/databases/{dbname}/tables/{tablename}/entries/{entryid}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		dbname := params["dbname"]
		tablename := params["tablename"]
		entryid := params["entryid"]
		deleter := &entries.EntryDeleter{DataDir: "./data", Logger: &logger.NoopLogger{}}
		err := deleter.DeleteEntry(dbname, tablename, entryid)
		var resp *response.APIResponse
		if err != nil {
			if err.Error() == "Eintrag nicht gefunden" {
				resp = response.WriteErrorInternal(404, "ERR_ENTRY_NOT_FOUND", err.Error(), "")
			} else {
				resp = response.WriteErrorInternal(400, "ERR_DELETE_ENTRY", err.Error(), "")
			}
		} else {
			resp = response.WriteSuccess(map[string]interface{}{
				"entryId": entryid,
				"deleted_at": time.Now().UTC().Format(time.RFC3339),
			}, "")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
	})

	// QueryHandler-Instanz für Query-Endpunkte
	queryHandler := &query.QueryHandler{DataDir: "./data", Logger: &logger.NoopLogger{}}

	// Query-Endpunkt: POST /api/databases/{dbname}/tables/{tablename}/query
	pr.HandleFunc("POST", "/api/databases/{dbname}/tables/{tablename}/query", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		dbName := params["dbname"]
		tableName := params["tablename"]
		resp := queryHandler.TableQueryHandler(dbName, tableName, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
	})
	// Query-Endpunkt: POST /api/{dbname}/query
	pr.HandleFunc("POST", "/api/{dbname}/query", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		dbName := params["dbname"]
		resp := queryHandler.QueryHandler(dbName, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
	})

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
