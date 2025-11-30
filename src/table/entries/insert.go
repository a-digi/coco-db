package entries

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
				var keyParts []string
				for _, f := range idxMeta.Fields {
					v, ok := entry[f]
					if !ok {
						keyParts = nil
						break
					}
					keyParts = append(keyParts, fmt.Sprint(v))
				}
				if keyParts == nil {
					continue
				}
				idxKey := strings.Join(keyParts, "|")
				idxPath := filepath.Join(tableDir, "indexes", "index_"+idxMeta.Name+".json")
				writeIndexFileAtomic(idxPath, func(idxObj map[string][]string) {
					idxObj[idxKey] = append(idxObj[idxKey], entryId)
				})
			}
		}
	}

	// 6. Indexdateien pro Feld in entries/<feldname>.json pflegen
	if err == nil {
		var meta fields.TableMeta
		if err := json.Unmarshal(metaFile, &meta); err == nil {
			indexedFields := map[string]struct{}{}
			for _, idxMeta := range meta.Indexes {
				for _, f := range idxMeta.Fields {
					indexedFields[f] = struct{}{}
				}
			}
			for f := range indexedFields {
				v, ok := entry[f]
				if !ok {
					continue
				}
				fieldIdxPath := filepath.Join(entriesDir, f+".json")
				writeIndexFileAtomic(fieldIdxPath, func(fieldIdx map[string][]string) {
					key := fmt.Sprint(v)
					fieldIdx[key] = append(fieldIdx[key], entryId)
				})
			}
		}
	}

	// 7. TotalEntries in meta.json erhöhen
	if metaFile, err := os.ReadFile(metaPath); err == nil {
		var meta fields.TableMeta
		if err := json.Unmarshal(metaFile, &meta); err == nil {
			meta.TotalEntries++
			if f, err := os.Create(metaPath); err == nil {
				_ = json.NewEncoder(f).Encode(meta)
				f.Close()
			}
		}
	}

	// Rückgabe: entryId und vollständiger Eintrag im Data-Objekt
	respData := map[string]interface{}{"entryId": entryId}
	for k, v := range entry {
		respData[k] = v
	}
	return &response.APIResponse{
		HttpCode:      201,
		Success:       true,
		Data:          respData,
		ExecutionTime: time.Since(start).String(),
	}
}

var (
	indexFileLocks   = make(map[string]*sync.Mutex)
	indexFileLocksMu sync.Mutex
)

// getIndexFileLock gibt einen Mutex für den gegebenen Dateipfad zurück (pro Datei eindeutig)
func getIndexFileLock(path string) *sync.Mutex {
	indexFileLocksMu.Lock()
	defer indexFileLocksMu.Unlock()
	m, ok := indexFileLocks[path]
	if !ok {
		m = &sync.Mutex{}
		indexFileLocks[path] = m
	}
	return m
}

// writeIndexFileAtomic führt Lesen, Modifizieren und Schreiben atomar unter Lock aus
func writeIndexFileAtomic(path string, updateFn func(map[string][]string)) {
	lock := getIndexFileLock(path)
	lock.Lock()
	defer lock.Unlock()
	idxObj := map[string][]string{}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &idxObj)
	}
	updateFn(idxObj)
	f, _ := os.Create(path)
	_ = json.NewEncoder(f).Encode(idxObj)
	f.Close()
}
