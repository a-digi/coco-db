package entries

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/a-digi/coco-db/src/table/fields/validate"
	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/index"
)

type EntryEditor struct {
	DataDir string
	Logger  logger.Logger
}

// EditEntry aktualisiert einen bestehenden Eintrag, versioniert die alte Version und speichert die neue als aktuelle Version.
func (ee *EntryEditor) EditEntry(dbName, tableName, entryId string, entry map[string]interface{}) error {
	if resp := validate.ValidateNoIDField(entry); resp != nil {
		return fmt.Errorf("%v", resp.Error.Message)
	}

	entryDir := filepath.Join(ee.DataDir, dbName, tableName, "entries", entryId)
	entryPath := filepath.Join(entryDir, entryId+".json")
	versioningObj := Versioning{Dir: entryDir}

	// Prüfen, ob der Eintrag existiert
	if _, err := os.Stat(entryPath); os.IsNotExist(err) {
		return fmt.Errorf("Eintrag nicht gefunden")
	}

	// 1. Versionierung der aktuellen Version
	versions, err := versioningObj.LoadVersions()
	if err != nil {
		ee.Logger.Error(fmt.Sprintf("Fehler beim Laden der Versionen: %v", err))
		return err
	}
	var currentVersion int
	if len(versions) > 0 {
		currentVersion = versions[len(versions)-1].VersionNumber
		// Alte Version archivieren, aber Version nicht erhöhen!
		if _, err := versioningObj.VersionActiveEntry(); err != nil {
			ee.Logger.Error(fmt.Sprintf("Fehler beim Versionieren des aktiven Eintrags: %v", err))
			return err
		}
	} else {
		return fmt.Errorf("Keine aktive Version vorhanden")
	}

	// id-Feld manuell setzen
	entry["id"] = entryId
	// 2. Neue Version speichern
	f, err := os.Create(entryPath)
	if err != nil {
		ee.Logger.Error(fmt.Sprintf("Fehler beim Anlegen der Eintragsdatei: %v", err))
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		ee.Logger.Error(fmt.Sprintf("Fehler beim Schreiben des Eintrags: %v", err))
		return err
	}

	// 3. Neue Version als aktuell eintragen
	if err := versioningObj.AddNewVersion(entryId, currentVersion+1); err != nil {
		ee.Logger.Error(fmt.Sprintf("Fehler beim Hinzufügen der neuen Version: %v", err))
		return err
	}

	// 4. Indexaktualisierung nach Edit
	tableDir := filepath.Join(ee.DataDir, dbName, tableName)
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
					// Lade alten Eintrag (vor Edit)
					oldEntryPath := filepath.Join(entryDir, entryId+".json")
					oldData, err := os.ReadFile(oldEntryPath)
					var oldEntry map[string]interface{}
					if err == nil {
						_ = json.Unmarshal(oldData, &oldEntry)
					}
					// Index laden
					idxPath := filepath.Join(tableDir, "indexes", "index_"+idxName+".json")
					var idxObj map[string][]string
					idxObj = map[string][]string{}
					if idxData, err := os.ReadFile(idxPath); err == nil {
						_ = json.Unmarshal(idxData, &idxObj)
					}
					// Alten Wert entfernen
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
						}
					}
					// Neuen Wert einfügen
					newKey, ok := entry[idxField].(string)
					if ok {
						idxObj[newKey] = append(idxObj[newKey], entryId)
						// RAM-Index vollständig aus Datei laden und ersetzen
						reg := index.GetRegistry()
						idxData, err := os.ReadFile(idxPath)
						var idxObjDisk map[string][]string
						if err == nil {
							_ = json.Unmarshal(idxData, &idxObjDisk)
							if ids, ok := idxObjDisk[newKey]; ok {
								reg.UpdateIndexInMemory(dbName, tableName, idxName, newKey, ids, "update")
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

	ee.Logger.Info(fmt.Sprintf("Eintrag erfolgreich aktualisiert: %s (Version %d)", entryPath, currentVersion+1))
	return nil
}
