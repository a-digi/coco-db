package entries_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	entries "github.com/a-digi/coco-db/src/table/entries"
)

func TestInsertEntry_APIResponse(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entrycreator_test")
	if err != nil {
		t.Fatalf("TempDir error: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbName := "testdb"
	tableName := "testtable"
	entry := map[string]interface{}{"foo": "bar", "num": 42}

	resp := entries.InsertEntry(dbName, tableName, entry, tmpDir, &logger.NoopLogger{})
	if resp == nil {
		t.Fatalf("InsertEntry gibt nil zurück")
	}
	if resp.HttpCode != 201 {
		t.Errorf("Erwartet HttpCode 201, erhalten: %d", resp.HttpCode)
	}
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data ist nicht vom Typ map[string]interface{}")
	}
	entryId, ok := data["entryId"].(string)
	if !ok || entryId == "" {
		t.Errorf("entryId fehlt oder ist leer: %+v", data)
	}
	// Prüfe, ob Datei existiert
	entriesDir := filepath.Join(tmpDir, dbName, tableName, "entries")
	entryDir := filepath.Join(entriesDir, entryId)
	entryFile := filepath.Join(entryDir, entryId+".json")
	if _, err := os.Stat(entryFile); err != nil {
		t.Errorf("Entry-Datei nicht gefunden: %v", err)
	}
	// Prüfe version.json
	versionFile := filepath.Join(entryDir, "version.json")
	fileData, err := os.ReadFile(versionFile)
	if err != nil {
		t.Errorf("version.json nicht gefunden: %v", err)
	}
	var versions []entries.VersionEntry
	if err := json.Unmarshal(fileData, &versions); err != nil {
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
