package entries_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/table/entries"
)

type testLogger struct{}
func (l *testLogger) Info(v ...interface{})  {}
func (l *testLogger) Error(v ...interface{}) {}
func (l *testLogger) Log(v ...interface{}) {}
func (l *testLogger) Debug(v ...interface{}) {}
func (l *testLogger) Notice(v ...interface{}) {}
func (l *testLogger) Warning(v ...interface{}) {}
func (l *testLogger) Critical(v ...interface{}) {}
func (l *testLogger) Alert(v ...interface{}) {}
func (l *testLogger) Emergency(v ...interface{}) {}

func setupTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_edit_entry_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownTestDir(dir string) {
	_ = os.RemoveAll(dir)
}

func createTestEntry(t *testing.T, dataDir, dbName, tableName, entryId string, entry map[string]interface{}) string {
	entryDir := filepath.Join(dataDir, dbName, tableName, "entries", entryId)
	if err := os.MkdirAll(entryDir, 0755); err != nil {
		t.Fatalf("Fehler beim Anlegen des entry-Ordners: %v", err)
	}
	entryPath := filepath.Join(entryDir, entryId+".json")
	f, err := os.Create(entryPath)
	if err != nil {
		t.Fatalf("Fehler beim Anlegen der Eintragsdatei: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		t.Fatalf("Fehler beim Schreiben des Eintrags: %v", err)
	}
	// Versionierung initialisieren
	versioningObj := entries.Versioning{Dir: entryDir}
	if err := versioningObj.AddNewVersion(entryId, 1); err != nil {
		t.Fatalf("Fehler beim Initialisieren der Versionierung: %v", err)
	}
	return entryPath
}

func TestEditEntry_Success(t *testing.T) {
	testDir := setupTestDir(t)
	defer teardownTestDir(testDir)
	logger := &testLogger{}
	dbName := "testdb"
	tableName := "users"
	entryId := "test-entry-1"
	origEntry := map[string]interface{}{"name": "Max", "age": 30}
	createTestEntry(t, testDir, dbName, tableName, entryId, origEntry)

	editor := &entries.EntryEditor{DataDir: testDir, Logger: logger}
	newEntry := map[string]interface{}{"name": "Max Mustermann", "age": 31}
	err := editor.EditEntry(dbName, tableName, entryId, newEntry)
	if err != nil {
		t.Fatalf("EditEntry Fehler: %v", err)
	}
	// Prüfen, ob die Datei überschrieben wurde
	entryPath := filepath.Join(testDir, dbName, tableName, "entries", entryId, entryId+".json")
	content, err := os.ReadFile(entryPath)
	if err != nil {
		t.Fatalf("Eintrag nicht lesbar: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(content, &got); err != nil {
		t.Fatalf("Eintrag nicht parsebar: %v", err)
	}
	if got["name"] != "Max Mustermann" || int(got["age"].(float64)) != 31 {
		t.Errorf("Eintrag nicht korrekt aktualisiert: %+v", got)
	}
	// Prüfen, ob die Versionierung hochgezählt wurde
	versioningObj := entries.Versioning{Dir: filepath.Join(testDir, dbName, tableName, "entries", entryId)}
	versions, err := versioningObj.LoadVersions()
	if err != nil {
		t.Fatalf("Fehler beim Laden der Versionen: %v", err)
	}
	if len(versions) != 2 || versions[1].VersionNumber != 2 {
		t.Errorf("Versionierung nicht korrekt: %+v", versions)
	}
}

func TestEditEntry_NotFound(t *testing.T) {
	testDir := setupTestDir(t)
	defer teardownTestDir(testDir)
	logger := &testLogger{}
	editor := &entries.EntryEditor{DataDir: testDir, Logger: logger}
	err := editor.EditEntry("testdb", "users", "not-exist", map[string]interface{}{"name": "Test"})
	if err == nil || err.Error() != "Eintrag nicht gefunden" {
		t.Errorf("Nicht vorhandener Eintrag nicht korrekt erkannt: %v", err)
	}
}

func TestEditEntry_InvalidField(t *testing.T) {
	testDir := setupTestDir(t)
	defer teardownTestDir(testDir)
	logger := &testLogger{}
	dbName := "testdb"
	tableName := "users"
	entryId := "test-entry-2"
	origEntry := map[string]interface{}{"name": "Anna"}
	createTestEntry(t, testDir, dbName, tableName, entryId, origEntry)

	editor := &entries.EntryEditor{DataDir: testDir, Logger: logger}
	invalidEntry := map[string]interface{}{"id": "should-not-be-here", "name": "Anna"}
	err := editor.EditEntry(dbName, tableName, entryId, invalidEntry)
	if err == nil || err.Error() == "" {
		t.Errorf("Ungültiges Feld nicht korrekt erkannt: %v", err)
	}
}

func TestEditEntry_FuncVariant(t *testing.T) {
	testDir := setupTestDir(t)
	defer teardownTestDir(testDir)
	dbName := "testdb"
	tableName := "users"
	entryId := "test-entry-3"
	origEntry := map[string]interface{}{"name": "Lisa"}
	createTestEntry(t, testDir, dbName, tableName, entryId, origEntry)

	newEntry := map[string]interface{}{"name": "Lisa Müller"}
	resp := entries.EditEntry(dbName, tableName, entryId, newEntry, testDir, &testLogger{})
	if resp == nil || !resp.Success {
		t.Fatalf("EditEntry (Funktion) fehlgeschlagen: %+v", resp)
	}
	// Prüfen, ob die Datei überschrieben wurde
	entryPath := filepath.Join(testDir, dbName, tableName, "entries", entryId, entryId+".json")
	content, err := os.ReadFile(entryPath)
	if err != nil {
		t.Fatalf("Eintrag nicht lesbar: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(content, &got); err != nil {
		t.Fatalf("Eintrag nicht parsebar: %v", err)
	}
	if got["name"] != "Lisa Müller" {
		t.Errorf("Eintrag nicht korrekt aktualisiert: %+v", got)
	}
}
