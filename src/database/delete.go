package database

import (
	"os"
	"path/filepath"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

type DatabaseDelete struct {
	DataDir string
	Logger  logger.Logger
}

func (dd *DatabaseDelete) HandleDeleteDatabase(name string) *response.APIResponse {
	if len(name) == 0 || !DbNamePattern.MatchString(name) {
		return response.WriteError(nil, 0, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
	}
	dbPath := filepath.Join(dd.DataDir, name)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return response.WriteError(nil, 0, "ERR_DB_NOT_FOUND", "Datenbank nicht gefunden", "")
	}
	if err := os.RemoveAll(dbPath); err != nil {
		return response.WriteError(nil, 0, "ERR_IO", err.Error(), "")
	}
	return response.WriteSuccess(nil, map[string]string{"name": name}, "")
}
