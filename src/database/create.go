package database

import (
	"os"
	"path/filepath"
    "github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

type DatabaseCreate struct {
	DataDir     string
	Logger      logger.Logger
}

func (dc *DatabaseCreate) HandleCreateDatabase(dbName string) *response.APIResponse {

	if !DbNamePattern.MatchString(dbName) {
		return response.WriteError(nil, 0, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
	}

	dbPath := filepath.Join(dc.DataDir, dbName)
	if _, err := os.Stat(dbPath); err == nil {
		return response.WriteError(nil, 0, "ERR_DB_EXISTS", "Datenbank existiert bereits", "")
	}

	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return response.WriteError(nil, 0, "ERR_IO", err.Error(), "")
	}

	return response.WriteSuccess(nil, map[string]string{"name": dbName}, "")
}
