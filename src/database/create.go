package database

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

type DatabaseCreate struct {
	DataDir     string
	Logger      logger.Logger
}

func (dc *DatabaseCreate) HandleCreateDatabase(dbName string) *response.APIResponse {
	start := time.Now()

	if !DbNamePattern.MatchString(dbName) {
		return response.WriteError(http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
	}

	dbPath := filepath.Join(dc.DataDir, dbName)
	if _, err := os.Stat(dbPath); err == nil {
		return response.WriteError(http.StatusConflict, "ERR_DB_EXISTS", "Datenbank existiert bereits", "")
	}

	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return response.WriteError(http.StatusInternalServerError, "ERR_IO", err.Error(), "")
	}

	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]string{"name": dbName}, execTime)
}
