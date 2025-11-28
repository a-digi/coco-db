package table

import (
	"net/http"
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
func (lh *ListTablesHandler) HandleListTables(w http.ResponseWriter, dbname string) {
	if dbname == "" {
		response.WriteError(w, http.StatusBadRequest, "ERR_DB_NAME_MISSING", "Datenbankname fehlt", "")
		return
	}
	dbDir := filepath.Join(lh.DataDir, dbname)
	dirs, err := os.ReadDir(dbDir)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "ERR_DB_NOT_FOUND", "Datenbank nicht gefunden", "")
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
	response.WriteSuccess(w, map[string]interface{}{ "tables": tables }, "Tabellen erfolgreich aufgelistet")
}
