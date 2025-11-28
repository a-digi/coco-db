package database

import (
	"os"
	"path/filepath"

	"github.com/a-digi/coco-db/src/response"
)

type UpdateDatabaseRequest struct {
	NewName string `json:"newName"`
}

type DatabaseUpdate struct {
	SataDir string
}

func (du *DatabaseUpdate) HandleUpdateDatabase(oldName, newName string) *response.APIResponse {
	if len(oldName) == 0 || !DbNamePattern.MatchString(oldName) {
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_INVALID_NAME",
				Message: "Ungültiger alter Datenbankname",
			},
		}
	}

	if len(newName) == 0 {
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_INVALID_JSON",
				Message: "newName fehlt",
			},
		}
	}

	if !DbNamePattern.MatchString(newName) {
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_INVALID_NAME",
				Message: "Ungültiger neuer Datenbankname",
			},
		}
	}

	if newName == oldName {
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_SAME_NAME",
				Message: "Neuer Name ist identisch mit altem Namen",
			},
		}
	}

	oldPath := filepath.Join(du.SataDir, oldName)
	newPath := filepath.Join(du.SataDir, newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_NOT_FOUND",
				Message: "Alte Datenbank nicht gefunden",
			},
		}
	}

	if _, err := os.Stat(newPath); err == nil {
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_DB_EXISTS",
				Message: "Ziel-Datenbankname existiert bereits",
			},
		}
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return &response.APIResponse{
			Success: false,
			Error: &response.APIError{
				Code:    "ERR_IO",
				Message: err.Error(),
			},
		}
	}
	return &response.APIResponse{
		Success: true,
		Data:    map[string]string{"oldName": oldName, "newName": newName},
	}
}
