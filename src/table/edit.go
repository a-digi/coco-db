package table

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/logger"
	"time"
)

// TableUpdate holds the metadata for updating a table
type TableUpdate struct {
	DataDir string
	Logger  logger.Logger
}

// HandleEditTable verarbeitet das Bearbeiten der Metadaten einer Tabelle (PUT /api/databases/{dbname}/tables/{tname})
func (tu *TableUpdate) HandleEditTable(w http.ResponseWriter, r *http.Request, dbname, tableName string) {
	start := time.Now()
	var meta TableMeta
	if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_JSON", "Ungültiges JSON: "+err.Error(), execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	if err := validateTableName(tableName); err != nil {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_TABLE_INVALID_NAME", err.Error(), execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	if err := validateFields(meta.Fields); err != nil {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusBadRequest, "ERR_FIELD_INVALID", err.Error(), execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	dbDir := filepath.Join(tu.DataDir, dbname)
	tableDir := filepath.Join(dbDir, tableName)
	metaPath := filepath.Join(tableDir, "meta.json")
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusNotFound, "ERR_TABLE_NOT_FOUND", "Tabelle nicht gefunden", execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	if meta.SchemaVersion == 0 {
		meta.SchemaVersion = 1
	}
	f, err := os.Create(metaPath)
	if err != nil {
		execTime := time.Since(start).String()
		td := fmt.Sprintf("Fehler beim Öffnen meta.json: %v", err)
		resp := response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", td, execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(meta); err != nil {
		execTime := time.Since(start).String()
		resp := response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Schreiben der meta.json: "+err.Error(), execTime)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.HttpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	tu.Logger.Info(fmt.Sprintf("[TABLE_EDIT] Tabelle %s/%s erfolgreich bearbeitet", dbname, tableName))
	execTime := time.Since(start).String()
	resp := response.WriteSuccess(meta, execTime)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.HttpCode)
	_ = json.NewEncoder(w).Encode(resp)
}
