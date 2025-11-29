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
func (tu *TableUpdate) HandleEditTable(dbname, tableName string, meta TableMeta) *response.APIResponse {
	start := time.Now()
	if err := validateTableName(tableName); err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_TABLE_INVALID_NAME", err.Error(), execTime)
	}
	if err := validateFields(meta.Fields); err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_FIELD_INVALID", err.Error(), execTime)
	}
	dbDir := filepath.Join(tu.DataDir, dbname)
	tableDir := filepath.Join(dbDir, tableName)
	metaPath := filepath.Join(tableDir, "meta.json")
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusNotFound, "ERR_TABLE_NOT_FOUND", "Tabelle nicht gefunden", execTime)
	}
	if meta.SchemaVersion == 0 {
		meta.SchemaVersion = 1
	}
	f, err := os.Create(metaPath)
	if err != nil {
		execTime := time.Since(start).String()
		td := fmt.Sprintf("Fehler beim Öffnen meta.json: %v", err)
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", td, execTime)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(meta); err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Schreiben der meta.json: "+err.Error(), execTime)
	}

	// tables.json der Datenbank aktualisieren
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	var tablesMeta []TableMeta
	if _, err := os.Stat(tablesJsonPath); err == nil {
		content, err := os.ReadFile(tablesJsonPath)
		if err == nil {
			_ = json.Unmarshal(content, &tablesMeta)
		}
	}
	// Ersetze die Metadaten der bearbeiteten Tabelle
	tableUpdated := false
	for i, t := range tablesMeta {
		if t.TableName == meta.TableName {
			tablesMeta[i] = meta
			tableUpdated = true
			break
		}
	}
	if !tableUpdated {
		tablesMeta = append(tablesMeta, meta)
	}
	f2, err := os.Create(tablesJsonPath)
	if err == nil {
		enc2 := json.NewEncoder(f2)
		enc2.SetIndent("", "  ")
		_ = enc2.Encode(tablesMeta)
		f2.Close()
	}

	tu.Logger.Info(fmt.Sprintf("[TABLE_EDIT] Tabelle %s/%s erfolgreich bearbeitet", dbname, tableName))
	execTime := time.Since(start).String()
	return response.WriteSuccess(meta, execTime)
}
