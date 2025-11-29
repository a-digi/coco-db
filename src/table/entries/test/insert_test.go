package entries_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	entries "github.com/a-digi/coco-db/src/table/entries"
)

func TestEntryCreator_InsertEntry_CreatesEntryAndVersion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entrycreator_test")
	if err != nil {
		t.Fatalf("TempDir error: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	ec := &entries.EntryCreator{
		DataDir: tmpDir,
		Logger:  &logger.NoopLogger{},
	}
	dbName := "testdb"
	tableName := "testtable"
	entry := map[string]interface{}{"foo": "bar", "num": 42}

	err = ec.InsertEntry(dbName, tableName, entry)
	if err != nil {
		t.Fatalf("InsertEntry error: %v", err)
	}

	// Suche das entries-Verzeichnis
	entriesDir := filepath.Join(tmpDir, dbName, tableName, "entries")
	files, err := os.ReadDir(entriesDir)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("Kein entryId-Ordner gefunden")
	}
	entryId := files[0].Name()
	entryDir := filepath.Join(entriesDir, entryId)
	entryFile := filepath.Join(entryDir, entryId+".json")
	if _, err := os.Stat(entryFile); err != nil {
		t.Errorf("Entry-Datei nicht gefunden: %v", err)
	}
	// Prüfe version.json
	versionFile := filepath.Join(entryDir, "version.json")
	data, err := os.ReadFile(versionFile)
	if err != nil {
		t.Errorf("version.json nicht gefunden: %v", err)
	}
	var versions []entries.VersionEntry
	if err := json.Unmarshal(data, &versions); err != nil {
		t.Errorf("Fehler beim Unmarshal von version.json: %v", err)
	}
	if len(versions) != 1 {
		t.Errorf("Erwartet 1 Versionseintrag, gefunden: %d", len(versions))
	}
	if versions[0].EntryID != entryId {
		t.Errorf("EntryID stimmt nicht: %s != %s", versions[0].EntryID, entryId)
	}
	if versions[0].VersionNumber != 1 {
		t.Errorf("VersionNumber stimmt nicht: %d", versions[0].VersionNumber)
	}
	if versions[0].VersionedAt != "" {
		t.Errorf("VersionedAt sollte leer sein, ist aber: %s", versions[0].VersionedAt)
	}
}
