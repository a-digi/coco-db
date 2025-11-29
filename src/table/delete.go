package table

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

// TableDelete kapselt die Abhängigkeiten für das Löschen von Tabellen
// DataDir: Basisverzeichnis der Datenbank, Logger: Logging-Instanz
//
type TableDelete struct {
	DataDir string
	Logger  logger.Logger
}

// HandleDeleteTable verarbeitet das Löschen einer Tabelle (DELETE /api/databases/{dbname}/tables/{tname})
func (td *TableDelete) HandleDeleteTable(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	dbname := r.URL.Query().Get("db")
	tableName := r.URL.Query().Get("table")
	if dbname == "" || tableName == "" {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_PARAM_MISSING", "Datenbank- oder Tabellenname fehlt", execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	dbDir := filepath.Join(td.DataDir, dbname)
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	var tablesMeta []TableMeta
	if _, err := os.Stat(tablesJsonPath); os.IsNotExist(err) {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusNotFound, "ERR_TABLES_NOT_FOUND", "tables.json nicht gefunden", execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	content, err := os.ReadFile(tablesJsonPath)
	if err != nil {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", err.Error(), execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	if err := json.Unmarshal(content, &tablesMeta); err != nil {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusInternalServerError, "ERR_JSON", err.Error(), execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	found := false
	newTables := make([]TableMeta, 0, len(tablesMeta))
	for _, meta := range tablesMeta {
		if strings.EqualFold(meta.TableName, tableName) {
			found = true
			continue // nicht übernehmen
		}
		newTables = append(newTables, meta)
	}
	if !found {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusNotFound, "ERR_TABLE_NOT_FOUND", "Tabelle nicht gefunden", execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	// Schreibe neue tables.json
	f, err := os.Create(tablesJsonPath)
	if err == nil {
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		_ = enc.Encode(newTables)
		f.Close()
	}
	// Lösche Tabellenverzeichnis
	tableDir := filepath.Join(dbDir, tableName)
	_ = os.RemoveAll(tableDir)
	td.Logger.Info("[TABLE_DELETE] Tabelle gelöscht:", dbname, tableName)
	execTime := time.Since(start).String()
	resp := response.WriteSuccess(map[string]string{"table": tableName}, execTime)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.HttpCode)
	_ = json.NewEncoder(w).Encode(resp)
}
