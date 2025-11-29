package entries

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/a-digi/coco-db/src/logger"
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

	// 1. Hard Delete: Lösche das gesamte Verzeichnis des Eintrags
	if err := versioningObj.HardDelete(); err != nil {
		ed.Logger.Error(fmt.Sprintf("Fehler beim Löschen des Eintrags: %v", err))
		return err
	}

	ed.Logger.Info(fmt.Sprintf("Eintrag erfolgreich gelöscht (hard): %s", entryPath))
	return nil
}
