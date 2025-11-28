package database

import (
	"encoding/json"
	"io/ioutil"
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

	// tables.json aktualisieren
	tablesPath := filepath.Join(dd.DataDir, "tables.json")
	var dbs []DatabaseMeta
	if _, err := os.Stat(tablesPath); err == nil {
		content, err := ioutil.ReadFile(tablesPath)
		if err == nil {
			_ = json.Unmarshal(content, &dbs)
		}
	}
	// Filtere die gelöschte DB heraus
	newDbs := make([]DatabaseMeta, 0, len(dbs))
	for _, entry := range dbs {
		if entry.Name != name {
			newDbs = append(newDbs, entry)
		}
	}
	f, err := os.Create(tablesPath)
	if err == nil {
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		_ = enc.Encode(newDbs)
		f.Close()
	}

	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]string{"name": name}, execTime)
}
