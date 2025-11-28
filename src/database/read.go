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

// DatabaseMeta ist bereits in create.go definiert und wird hier nicht erneut benötigt.

type DatabaseList struct {
	DataDir string
	Logger  logger.Logger
}

func (dl *DatabaseList) HandleListDatabases() *response.APIResponse {
	start := time.Now()
	dbJsonPath := filepath.Join(dl.DataDir, "db.json")
	var dbs []DatabaseMeta
	if _, err := os.Stat(dbJsonPath); os.IsNotExist(err) {
		execTime := time.Since(start).String()
		return response.WriteSuccess(map[string]interface{}{"databases": []DatabaseMeta{}}, execTime)
	}
	content, err := ioutil.ReadFile(dbJsonPath)
	if err != nil {
		execTime := time.Since(start).String()
		dl.Logger.Error("[DB_LIST] Fehler beim Lesen von db.json:", err)
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", err.Error(), execTime)
	}
	if err := json.Unmarshal(content, &dbs); err != nil {
		execTime := time.Since(start).String()
		dl.Logger.Error("[DB_LIST] Fehler beim Parsen von db.json:", err)
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_JSON", err.Error(), execTime)
	}
	dl.Logger.Info("[DB_LIST] Datenbanken aufgelistet:", dbs)
	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]interface{}{"databases": dbs}, execTime)
}
