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

// TableListMeta kapselt die Abhängigkeiten für das Listen von Tabellen-Metadaten
// DataDir: Basisverzeichnis der Datenbank, Logger: Logging-Instanz
//
type TableListMeta struct {
	DataDir string
	Logger  logger.Logger
}

// HandleGetTable verarbeitet das Abrufen der Metadaten einer Tabelle (GET /api/databases/{dbname}/tables/{tname})
func (tlm *TableListMeta) HandleGetTable(w http.ResponseWriter, r *http.Request) {
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
	dbDir := filepath.Join(tlm.DataDir, dbname)
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
	for _, meta := range tablesMeta {
		if strings.EqualFold(meta.TableName, tableName) {
			execTime := time.Since(start).String()
			resp := response.WriteSuccess(meta, execTime)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
	}
	execTime := time.Since(start).String()
	resp := response.WriteErrorInternal(http.StatusNotFound, "ERR_TABLE_NOT_FOUND", "Tabelle nicht gefunden", execTime)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.HttpCode)
	_ = json.NewEncoder(w).Encode(resp)
}
