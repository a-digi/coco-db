package database

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

type DatabaseDelete struct {
	DataDir string
	Logger  logger.Logger
}

func (dd *DatabaseDelete) HandleDeleteDatabase(name string) *response.APIResponse {
	start := time.Now()
	if len(name) == 0 || !DbNamePattern.MatchString(name) {
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
	}
	dbPath := filepath.Join(dd.DataDir, name)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return response.WriteErrorInternal(http.StatusNotFound, "ERR_DB_NOT_FOUND", "Datenbank nicht gefunden", "")
	}
	if err := os.RemoveAll(dbPath); err != nil {
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", err.Error(), "")
	}
	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]string{"name": name}, execTime)
}
