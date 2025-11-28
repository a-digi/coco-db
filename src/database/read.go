package database

import (
	"io/ioutil"
	"net/http"
	"github.com/a-digi/coco-db/src/response"
	"os"
)

func HandleListDatabases(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	return handleListDatabases(w, r, dataDir)
}

func handleListDatabases(w http.ResponseWriter, r *http.Request, dataDir string) *response.APIResponse {
	entries, err := ioutil.ReadDir(dataDir)
	if err != nil {
		// Wenn das Verzeichnis nicht existiert, gib eine leere Liste zurück
		if os.IsNotExist(err) {
			return response.WriteSuccess(w, map[string]interface{}{ "databases": []string{} }, "")
		}
		return response.WriteError(w, http.StatusInternalServerError, "ERR_IO", err.Error(), "")
	}
	dbs := []string{}
	for _, entry := range entries {
		if entry.IsDir() && DbNamePattern.MatchString(entry.Name()) {
			dbs = append(dbs, entry.Name())
		}
	}
	return response.WriteSuccess(w, map[string]interface{}{ "databases": dbs }, "")
}
