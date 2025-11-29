package entries

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	fields "github.com/a-digi/coco-db/src/table/fields"
	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

type EntryCreator struct {
	DataDir string
	Logger  logger.Logger
}

// InsertEntry speichert einen neuen Eintrag, legt die Versionierung an und loggt das Ergebnis.
func (ec *EntryCreator) InsertEntry(dbName, tableName string, entry map[string]interface{}) error {
	tableDir := filepath.Join(ec.DataDir, dbName, tableName)
	entriesDir := filepath.Join(tableDir, "entries")
	if err := os.MkdirAll(entriesDir, 0755); err != nil {
		ec.Logger.Error(fmt.Sprintf("Fehler beim Anlegen des entries-Verzeichnisses: %v", err))
		return fmt.Errorf("Fehler beim Anlegen des entries-Verzeichnisses: %v", err)
	}

	// 1. entryId generieren (UUIDv4)
	entryId, err := fields.NewUUIDv4()
	if err != nil {
		ec.Logger.Error(fmt.Sprintf("Fehler beim Generieren der UUID: %v", err))
		return fmt.Errorf("Fehler beim Generieren der UUID: %v", err)
	}
	entryDir := filepath.Join(entriesDir, entryId)
	if err := os.MkdirAll(entryDir, 0755); err != nil {
		ec.Logger.Error(fmt.Sprintf("Fehler beim Anlegen des entry-Ordners: %v", err))
		return fmt.Errorf("Fehler beim Anlegen des entry-Ordners: %v", err)
	}
	entryPath := filepath.Join(entryDir, entryId+".json")

	// 2. Versionierung vorbereiten
	versioningObj := Versioning{Dir: entryDir}
	var versionNumber int
	versions, err := versioningObj.LoadVersions()
	if err != nil {
		ec.Logger.Error(fmt.Sprintf("Fehler beim Laden der Versionen: %v", err))
		return fmt.Errorf("Fehler beim Laden der Versionen: %v", err)
	}
	if len(versions) > 0 {
		versionNumber, err = versioningObj.VersionActiveEntry()
		if err != nil {
			ec.Logger.Error(fmt.Sprintf("Fehler beim Versionieren des aktiven Eintrags: %v", err))
			return fmt.Errorf("Fehler beim Versionieren des aktiven Eintrags: %v", err)
		}
	} else {
		versionNumber = 1
	}

	// 3. Eintrag speichern
	f, err := os.Create(entryPath)
	if err != nil {
		ec.Logger.Error(fmt.Sprintf("Fehler beim Anlegen der Eintragsdatei: %v", err))
		return fmt.Errorf("Fehler beim Anlegen der Eintragsdatei: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		ec.Logger.Error(fmt.Sprintf("Fehler beim Schreiben des Eintrags: %v", err))
		return fmt.Errorf("Fehler beim Schreiben des Eintrags: %v", err)
	}

	// 4. Neue Version als aktuell eintragen
	if err := versioningObj.AddNewVersion(entryId, versionNumber); err != nil {
		ec.Logger.Error(fmt.Sprintf("Fehler beim Hinzufügen der neuen Version: %v", err))
		return fmt.Errorf("Fehler beim Hinzufügen der neuen Version: %v", err)
	}

	ec.Logger.Info(fmt.Sprintf("Eintrag erfolgreich gespeichert: %s (Version %d)", entryPath, versionNumber))
	return nil
}

// Funktionsbasierte Variante für direkten Aufruf
func InsertEntry(dbName, tableName string, entry map[string]interface{}, dataDir string, log logger.Logger) *response.APIResponse {
	tableDir := filepath.Join(dataDir, dbName, tableName)
	entriesDir := filepath.Join(tableDir, "entries")
	if err := os.MkdirAll(entriesDir, 0755); err != nil {
		log.Error(fmt.Sprintf("Fehler beim Anlegen des entries-Verzeichnisses: %v", err))
		return response.WriteErrorInternal(500, "ERR_CREATE_ENTRIES_DIR", err.Error(), "")
	}

	// 1. entryId generieren (UUIDv4)
	entryId, err := fields.NewUUIDv4()
	if err != nil {
		log.Error(fmt.Sprintf("Fehler beim Generieren der UUID: %v", err))
		return response.WriteErrorInternal(500, "ERR_UUID", err.Error(), "")
	}
	entryDir := filepath.Join(entriesDir, entryId)
	if err := os.MkdirAll(entryDir, 0755); err != nil {
		log.Error(fmt.Sprintf("Fehler beim Anlegen des entry-Ordners: %v", err))
		return response.WriteErrorInternal(500, "ERR_CREATE_ENTRY_DIR", err.Error(), "")
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
		return response.WriteErrorInternal(500, "ERR_CREATE_ENTRY_FILE", err.Error(), "")
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		log.Error(fmt.Sprintf("Fehler beim Schreiben des Eintrags: %v", err))
		return response.WriteErrorInternal(400, "ERR_WRITE_ENTRY", err.Error(), "")
	}

	// 4. Neue Version als aktuell eintragen
	if err := versioningObj.AddNewVersion(entryId, versionNumber); err != nil {
		log.Error(fmt.Sprintf("Fehler beim Hinzufügen der neuen Version: %v", err))
		return response.WriteErrorInternal(500, "ERR_ADD_VERSION", err.Error(), "")
	}

	log.Info(fmt.Sprintf("Eintrag erfolgreich gespeichert: %s (Version %d)", entryPath, versionNumber))
	return &response.APIResponse{
		HttpCode: 201,
		Data: map[string]interface{}{ "entryId": entryId },
	}
}
