package table

import (
	"encoding/json"
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
func (tlm *TableListMeta) HandleGetTable(dbname, tableName string) *response.APIResponse {
	start := time.Now()
	if dbname == "" || tableName == "" {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(400, "ERR_PARAM_MISSING", "Datenbank- oder Tabellenname fehlt", execTime)
	}
	dbDir := filepath.Join(tlm.DataDir, dbname)
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	var tablesMeta []TableMeta
	if _, err := os.Stat(tablesJsonPath); os.IsNotExist(err) {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(404, "ERR_TABLES_NOT_FOUND", "tables.json nicht gefunden", execTime)
	}
	content, err := os.ReadFile(tablesJsonPath)
	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(500, "ERR_IO", err.Error(), execTime)
	}
	if err := json.Unmarshal(content, &tablesMeta); err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(500, "ERR_JSON", err.Error(), execTime)
	}
	for _, meta := range tablesMeta {
		if strings.EqualFold(meta.TableName, tableName) {
			execTime := time.Since(start).String()
			return response.WriteSuccess(meta, execTime)
		}
	}
	execTime := time.Since(start).String()
	return response.WriteErrorInternal(404, "ERR_TABLE_NOT_FOUND", "Tabelle nicht gefunden", execTime)
}
