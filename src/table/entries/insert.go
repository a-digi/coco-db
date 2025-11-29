package entries

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	fields "github.com/a-digi/coco-db/src/table/fields"
	validate "github.com/a-digi/coco-db/src/table/fields/validate"
	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

// Funktionsbasierte Variante für direkten Aufruf
func InsertEntry(dbName, tableName string, entry map[string]interface{}, dataDir string, log logger.Logger) *response.APIResponse {
	start := time.Now()
	if resp := validate.ValidateNoIDField(entry); resp != nil {
		execTime := time.Since(start).String()
		resp.ExecutionTime = execTime
		return resp
	}
	tableDir := filepath.Join(dataDir, dbName, tableName)
	entriesDir := filepath.Join(tableDir, "entries")
	if err := os.MkdirAll(entriesDir, 0755); err != nil {
		log.Error(fmt.Sprintf("Fehler beim Anlegen des entries-Verzeichnisses: %v", err))
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(500, "ERR_CREATE_ENTRIES_DIR", err.Error(), execTime)
	}

	// 1. entryId generieren (UUIDv4)
	entryId, err := fields.NewUUIDv4()
	if err != nil {
		log.Error(fmt.Sprintf("Fehler beim Generieren der UUID: %v", err))
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(500, "ERR_UUID", err.Error(), execTime)
	}
	// id-Feld manuell setzen
	entry["id"] = entryId
	entryDir := filepath.Join(entriesDir, entryId)
	if err := os.MkdirAll(entryDir, 0755); err != nil {
		log.Error(fmt.Sprintf("Fehler beim Anlegen des entry-Ordners: %v", err))
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(500, "ERR_CREATE_ENTRY_DIR", err.Error(), execTime)
	}
	entryPath := filepath.Join(entryDir, entryId+".json")

	// 2. Versionierung vorbereiten
	versioningObj := Versioning{Dir: entryDir}
	var versionNumber int
	versions, err := versioningObj.LoadVersions()
	if err != nil {
		log.Error(fmt.Sprintf("Fehler beim Laden der Versionen: %v", err))
		return response.WriteErrorInternal(500, "ERR_LOAD_VERSIONS", err.Error(), "")
	}
	if len(versions) > 0 {
		versionNumber, err = versioningObj.VersionActiveEntry()
		if err != nil {
			log.Error(fmt.Sprintf("Fehler beim Versionieren des aktiven Eintrags: %v", err))
			return response.WriteErrorInternal(409, "ERR_VERSION_ACTIVE", err.Error(), "")
		}
	} else {
		versionNumber = 1
	}

	// 3. Eintrag speichern
	f, err := os.Create(entryPath)
	if err != nil {
		log.Error(fmt.Sprintf("Fehler beim Anlegen der Eintragsdatei: %v", err))
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(500, "ERR_CREATE_ENTRY_FILE", err.Error(), execTime)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		log.Error(fmt.Sprintf("Fehler beim Schreiben des Eintrags: %v", err))
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(400, "ERR_WRITE_ENTRY", err.Error(), execTime)
	}

	// 4. Neue Version als aktuell eintragen
	if err := versioningObj.AddNewVersion(entryId, versionNumber); err != nil {
		log.Error(fmt.Sprintf("Fehler beim Hinzufügen der neuen Version: %v", err))
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(500, "ERR_ADD_VERSION", err.Error(), execTime)
	}

	// 5. Indexaktualisierung nach Insert
	metaPath := filepath.Join(tableDir, "meta.json")
	metaFile, err := os.ReadFile(metaPath)
	if err == nil {
		var meta fields.TableMeta
		if err := json.Unmarshal(metaFile, &meta); err == nil {
			for _, idxMeta := range meta.Indexes {
				if len(idxMeta.Fields) == 1 {
					idxField := idxMeta.Fields[0]
					key, ok := entry[idxField]
					if ok {
						// Index laden oder neu anlegen
						idxPath := filepath.Join(tableDir, "index_"+idxMeta.Name+".json")
						var idxObj map[string][]string
						idxObj = map[string][]string{}
						if idxData, err := os.ReadFile(idxPath); err == nil {
							_ = json.Unmarshal(idxData, &idxObj)
						}
						k, ok := key.(string)
						if ok {
							idxObj[k] = append(idxObj[k], entryId)
							idxFile, _ := os.Create(idxPath)
							_ = json.NewEncoder(idxFile).Encode(idxObj)
							idxFile.Close()
						}
					}
				}
			}
		}
	}

	return &response.APIResponse{
		HttpCode:      201,
		Success:       true,
		Data:          entry,
		ExecutionTime: time.Since(start).String(),
	}
}
