package table

import (
	"os"
	"path/filepath"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/logger"
)

type ListTablesHandler struct {
	DataDir      string
	Logger       logger.Logger
	APIResponse  *response.APIResponse
}

// HandleListTables verarbeitet das Auflisten aller Tabellen (GET /api/databases/{dbname}/tables)
func (lh *ListTablesHandler) HandleListTables(dbname string) {
	if dbname == "" {
		lh.APIResponse = &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_NAME_MISSING",
				Message: "Datenbankname fehlt",
			},
		}
		return
	}
	dbDir := filepath.Join(lh.DataDir, dbname)
	dirs, err := os.ReadDir(dbDir)
	if err != nil {
		lh.APIResponse = &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_NOT_FOUND",
				Message: "Datenbank nicht gefunden",
			},
		}
		lh.Logger.Error("[TABLE_LIST] Datenbank nicht gefunden:", dbDir)
		return
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
	lh.APIResponse = &response.APIResponse{
		Success: true,
		Data: map[string]interface{}{ "tables": tables },
	}
}
