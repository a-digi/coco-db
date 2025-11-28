package table

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
	"encoding/json"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/logger"
)

type ListTablesHandler struct {
	DataDir string
	Logger  logger.Logger
}

// HandleListTables verarbeitet das Auflisten aller Tabellen (GET /api/databases/{dbname}/tables)
func (lh *ListTablesHandler) HandleListTables(dbname string) *response.APIResponse {
	start := time.Now()
	if dbname == "" {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_DB_NAME_MISSING", http.StatusText(http.StatusBadRequest)+": Datenbankname fehlt", execTime)
	}
	dbDir := filepath.Join(lh.DataDir, dbname)
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	var tablesMeta []TableMeta
	if _, err := os.Stat(tablesJsonPath); os.IsNotExist(err) {
		// Lege leeres Array an, falls Datei nicht existiert
		_ = os.WriteFile(tablesJsonPath, []byte("[]"), 0644)
		execTime := time.Since(start).String()
		return response.WriteSuccess(map[string]interface{}{ "tables": []string{} }, execTime)
	}
	content, err := os.ReadFile(tablesJsonPath)
	if err != nil {
		execTime := time.Since(start).String()
		lh.Logger.Error("[TABLE_LIST] Fehler beim Lesen von tables.json:", err)
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", err.Error(), execTime)
	}
	if err := json.Unmarshal(content, &tablesMeta); err != nil {
		execTime := time.Since(start).String()
		lh.Logger.Error("[TABLE_LIST] Fehler beim Parsen von tables.json:", err)
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_JSON", err.Error(), execTime)
	}
	tables := make([]string, 0, len(tablesMeta))
	for _, t := range tablesMeta {
		tables = append(tables, t.TableName)
	}
	lh.Logger.Info("[TABLE_LIST] Tabellen aufgelistet für DB:", dbname, tables)
	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]interface{}{ "tables": tables }, execTime)
}
