package table

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	dirs, err := os.ReadDir(dbDir)
	if err != nil {
		lh.Logger.Error("[TABLE_LIST] Datenbank nicht gefunden:", dbDir)
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusNotFound, "ERR_DB_NOT_FOUND", http.StatusText(http.StatusNotFound)+": Datenbank nicht gefunden", execTime)
	}
	tables := []string{}
	for _, entry := range dirs {
		if entry.IsDir() {
			metaPath := filepath.Join(dbDir, entry.Name(), "meta.json")
			if _, err := os.Stat(metaPath); err == nil {
				tables = append(tables, entry.Name())
			}
		}
	}
	lh.Logger.Info("[TABLE_LIST] Tabellen aufgelistet für DB:", dbname, tables)
	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]interface{}{ "tables": tables }, execTime)
}
