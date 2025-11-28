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

type DatabaseCreate struct {
	DataDir     string
	Logger      logger.Logger
}

// Struktur für einen Eintrag in tables.json
// (kann bei Bedarf erweitert werden)
type DatabaseMeta struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (dc *DatabaseCreate) HandleCreateDatabase(dbName string) *response.APIResponse {
	start := time.Now()

	if !DbNamePattern.MatchString(dbName) {
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
	}

	dbPath := filepath.Join(dc.DataDir, dbName)
	if _, err := os.Stat(dbPath); err == nil {
		return response.WriteErrorInternal(http.StatusConflict, "ERR_DB_EXISTS", "Datenbank existiert bereits", "")
	}

	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", err.Error(), "")
	}

	dbPath = filepath.Join(dc.DataDir, "db.json")
	var dbs []DatabaseMeta
	if _, err := os.Stat(dbPath); err == nil {
		content, err := ioutil.ReadFile(dbPath)
		if err == nil {
			_ = json.Unmarshal(content, &dbs)
		}
	}
	// Prüfe, ob der Name schon in tables.json existiert
	for _, entry := range dbs {
		if entry.Name == dbName {
			return response.WriteErrorInternal(http.StatusConflict, "ERR_DB_EXISTS", "Datenbank existiert bereits", "")
		}
	}
	dbs = append(dbs, DatabaseMeta{Name: dbName, CreatedAt: time.Now()})
	f, err := os.Create(dbPath)
	if err != nil {
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Schreiben von tables.json: "+err.Error(), "")
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(dbs); err != nil {
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Schreiben von tables.json: "+err.Error(), "")
	}

	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]string{"name": dbName}, execTime)
}
