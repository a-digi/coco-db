package table

import (
	"os"
	"path/filepath"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/logger"
)

type ListTablesHandler struct {
	DataDir string
	Logger  logger.Logger
}

// HandleListTables verarbeitet das Auflisten aller Tabellen (GET /api/databases/{dbname}/tables)
func (lh *ListTablesHandler) HandleListTables(dbname string) *response.APIResponse {
	if dbname == "" {
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_NAME_MISSING",
				Message: "Datenbankname fehlt",
			},
		}
	}
	dbDir := filepath.Join(lh.DataDir, dbname)
	dirs, err := os.ReadDir(dbDir)
	if err != nil {
		lh.Logger.Error("[TABLE_LIST] Datenbank nicht gefunden:", dbDir)
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_NOT_FOUND",
				Message: "Datenbank nicht gefunden",
			},
		}
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
	return &response.APIResponse{
		Success: true,
		Data: map[string]interface{}{ "tables": tables },
	}
}
