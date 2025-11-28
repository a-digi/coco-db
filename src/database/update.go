package database

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/a-digi/coco-db/src/response"
)

type UpdateDatabaseRequest struct {
	NewName string `json:"newName"`
}

type DatabaseUpdate struct {
	SataDir string
}

func (du *DatabaseUpdate) HandleUpdateDatabase(oldName, newName string) *response.APIResponse {
	start := time.Now()
	if len(oldName) == 0 || !DbNamePattern.MatchString(oldName) {
		execTime := time.Since(start).String()
		return response.WriteError(http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger alter Datenbankname", execTime)
	}

	if len(newName) == 0 {
		execTime := time.Since(start).String()
		return response.WriteError(http.StatusBadRequest, "ERR_DB_INVALID_JSON", "newName fehlt", execTime)
	}

	if !DbNamePattern.MatchString(newName) {
		execTime := time.Since(start).String()
		return response.WriteError(http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger neuer Datenbankname", execTime)
	}

	if newName == oldName {
		execTime := time.Since(start).String()
		return response.WriteError(http.StatusBadRequest, "ERR_DB_SAME_NAME", "Neuer Name ist identisch mit altem Namen", execTime)
	}

	oldPath := filepath.Join(du.SataDir, oldName)
	newPath := filepath.Join(du.SataDir, newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		execTime := time.Since(start).String()
		return response.WriteError(http.StatusNotFound, "ERR_DB_NOT_FOUND", "Alte Datenbank nicht gefunden", execTime)
	}

	if _, err := os.Stat(newPath); err == nil {
		execTime := time.Since(start).String()
		return response.WriteError(http.StatusConflict, "ERR_DB_EXISTS", "Ziel-Datenbankname existiert bereits", execTime)
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		execTime := time.Since(start).String()
		return response.WriteError(http.StatusInternalServerError, "ERR_IO", err.Error(), execTime)
	}
	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]string{"oldName": oldName, "newName": newName}, execTime)
}
