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
	APIResponse *response.APIResponse
}

func (dc *DatabaseCreate) HandleCreateDatabase(dbName string) *response.APIResponse {
	if !DbNamePattern.MatchString(dbName) {
		dc.APIResponse = &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_INVALID_NAME",
				Message: "Ungültiger Datenbankname",
			},
		}

		return dc.APIResponse
	}

	dbPath := filepath.Join(dc.DataDir, dbName)
	if _, err := os.Stat(dbPath); err == nil {
		dc.APIResponse = &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_EXISTS",
				Message: "Datenbank existiert bereits",
			},
		}
		return dc.APIResponse
	}

	if err := os.MkdirAll(dbPath, 0755); err != nil {
		dc.APIResponse = &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_IO",
				Message: err.Error(),
			},
		}
		return dc.APIResponse
	}

	dc.APIResponse = &response.APIResponse{
		Success: true,
		Data:    map[string]string{"name": dbName},
	}

	return dc.APIResponse
}
