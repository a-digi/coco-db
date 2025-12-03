package entries

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/index"
)

type EntryDeleter struct {
	DataDir string
	Logger  logger.Logger
}

// DeleteEntry markiert einen Eintrag als gelöscht (Soft Delete) und versioniert ihn.
func (ed *EntryDeleter) DeleteEntry(dbName, tableName, entryId string) error {
	entryDir := filepath.Join(ed.DataDir, dbName, tableName, "entries", entryId)
	entryPath := filepath.Join(entryDir, entryId+".json")
	versioningObj := Versioning{Dir: entryDir}

	// Prüfen, ob der Eintrag existiert
	if _, err := os.Stat(entryPath); os.IsNotExist(err) {
		return fmt.Errorf("Eintrag nicht gefunden")
	}

	// Fix: Lade den alten Eintrag vor dem Hard Delete
	var oldEntry map[string]interface{}
	if oldData, err := os.ReadFile(entryPath); err == nil {
		_ = json.Unmarshal(oldData, &oldEntry)
	}

	// 1. Hard Delete: Lösche das gesamte Verzeichnis des Eintrags
	if err := versioningObj.HardDelete(); err != nil {
		ed.Logger.Error(fmt.Sprintf("Fehler beim Löschen des Eintrags: %v", err))
		return err
	}

	// 2. Indexaktualisierung nach Delete
	tableDir := filepath.Join(ed.DataDir, dbName, tableName)
	metaPath := filepath.Join(tableDir, "meta.json")
	metaFile, err := os.ReadFile(metaPath)
	if err == nil {
		var meta map[string]interface{}
		if err := json.Unmarshal(metaFile, &meta); err == nil {
			indexes, ok := meta["indexes"].([]interface{})
			if ok {
				for _, idxRaw := range indexes {
					idxMeta, ok := idxRaw.(map[string]interface{})
					if !ok { continue }
					fieldsArr, ok := idxMeta["fields"].([]interface{})
					if !ok || len(fieldsArr) != 1 { continue }
					idxField, _ := fieldsArr[0].(string)
					idxName, _ := idxMeta["name"].(string)
					// Index laden
					idxPath := filepath.Join(tableDir, "indexes", "index_"+idxName+".json")
					var idxObj map[string][]string
					idxObj = map[string][]string{}
					if idxData, err := os.ReadFile(idxPath); err == nil {
						_ = json.Unmarshal(idxData, &idxObj)
					}
					// Wert entfernen
					if oldEntry != nil {
						oldKey, ok := oldEntry[idxField].(string)
						if ok {
							ids := idxObj[oldKey]
							newIds := []string{}
							for _, id := range ids {
								if id != entryId {
									newIds = append(newIds, id)
								}
							}
							if len(newIds) > 0 {
								idxObj[oldKey] = newIds
							} else {
								delete(idxObj, oldKey)
							}
							// RAM-Index vollständig aus Datei laden und ersetzen
							reg := index.GetRegistry()
							idxData, err := os.ReadFile(idxPath)
							var idxObjDisk map[string][]string
							if err == nil {
								_ = json.Unmarshal(idxData, &idxObjDisk)
								if ids, ok := idxObjDisk[oldKey]; ok && len(ids) > 0 {
									reg.UpdateIndexInMemory(dbName, tableName, idxName, oldKey, ids, "delete")
								} else {
									reg.UpdateIndexInMemory(dbName, tableName, idxName, oldKey, nil, "delete")
								}
							}
						}
					}
					// Index speichern
					idxFile, _ := os.Create(idxPath)
					_ = json.NewEncoder(idxFile).Encode(idxObj)
					idxFile.Close()
				}
			}
		}
	}

	// 3. TotalEntries in meta.json verringern
	metaFile, err = os.ReadFile(metaPath)
	if err == nil {
		var meta map[string]interface{}
		if err := json.Unmarshal(metaFile, &meta); err == nil {
			if te, ok := meta["totalEntries"].(float64); ok {
				if te > 0 {
					meta["totalEntries"] = te - 1
				}
			}
			if f, err := os.Create(metaPath); err == nil {
				_ = json.NewEncoder(f).Encode(meta)
				f.Close()
			}
		}
	}

	ed.Logger.Info(fmt.Sprintf("Eintrag erfolgreich gelöscht (hard): %s", entryPath))
	return nil
}
