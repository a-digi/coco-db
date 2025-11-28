package database

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

type DatabaseList struct {
	DataDir string
	Logger  logger.Logger
}

func (dl *DatabaseList) HandleListDatabases() *response.APIResponse {
	start := time.Now()
	entries, err := ioutil.ReadDir(dl.DataDir)
	if err != nil {
		execTime := time.Since(start).String()
		if os.IsNotExist(err) {
			return response.WriteSuccess(map[string]interface{}{"databases": []string{}}, execTime)
		}
		dl.Logger.Error("[DB_LIST] Fehler beim Lesen des Datenbankverzeichnisses:", err)
		return response.WriteError(http.StatusInternalServerError, "ERR_IO", err.Error(), execTime)
	}
	dbs := []string{}
	for _, entry := range entries {
		if entry.IsDir() && DbNamePattern.MatchString(entry.Name()) {
			dbs = append(dbs, entry.Name())
		}
	}
	dl.Logger.Info("[DB_LIST] Datenbanken aufgelistet:", dbs)
	execTime := time.Since(start).String()
	return response.WriteSuccess(map[string]interface{}{"databases": dbs}, execTime)
}
