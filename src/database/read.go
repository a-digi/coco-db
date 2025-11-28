package database

import (
	"io/ioutil"
	"os"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/logger"
)

type DatabaseList struct {
	DataDir string
	Logger  logger.Logger
}

func (dl *DatabaseList) HandleListDatabases() *response.APIResponse {
	entries, err := ioutil.ReadDir(dl.DataDir)

	if err != nil {
		if os.IsNotExist(err) {
			return response.WriteSuccess(nil, map[string]interface{}{ "databases": []string{} }, "")
		}
		dl.Logger.Error("[DB_LIST] Fehler beim Lesen des Datenbankverzeichnisses:", err)
		return response.WriteError(nil, 0, "ERR_IO", err.Error(), "")
	}
	dbs := []string{}
	for _, entry := range entries {
		if entry.IsDir() && DbNamePattern.MatchString(entry.Name()) {
			dbs = append(dbs, entry.Name())
		}
	}
	dl.Logger.Info("[DB_LIST] Datenbanken aufgelistet:", dbs)

	return response.WriteSuccess(nil, map[string]interface{}{ "databases": dbs }, "")
}
