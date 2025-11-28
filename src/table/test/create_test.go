package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/table"
)

func setupCreateTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_create_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownCreateTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestHandleCreateTable_Success(t *testing.T) {
	testDir := setupCreateTestDir(t)
	defer teardownCreateTestDir(testDir)
	creator := &table.TableCreator{
		DataDir: testDir,
		Logger:  &logger.NoopLogger{},
	}
	meta := table.TableMeta{
		TableName: "users",
		Fields:    []table.FieldMeta{{Name: "id", Type: "string"}},
	}
	resp := creator.HandleCreateTable("testdb", meta)
	if resp == nil || !resp.Success {
		t.Fatalf("Erwartet: Success true, erhalten: %+v", resp)
	}
	// Prüfe, ob meta.json existiert
	metaPath := filepath.Join(testDir, "testdb", "users", "meta.json")
	if _, err := os.Stat(metaPath); err != nil {
		t.Errorf("meta.json wurde nicht angelegt: %v", err)
	}
	// Prüfe, ob tables.json existiert und die Tabelle eingetragen ist
	tablesJsonPath := filepath.Join(testDir, "testdb", "tables.json")
	content, err := os.ReadFile(tablesJsonPath)
	if err != nil {
		t.Fatalf("tables.json wurde nicht angelegt: %v", err)
	}
	var tablesMeta []table.TableMeta
	if err := json.Unmarshal(content, &tablesMeta); err != nil {
		t.Fatalf("tables.json nicht parsebar: %v", err)
	}
	found := false
	for _, entry := range tablesMeta {
		if entry.TableName == "users" {
			found = true
		}
	}
	if !found {
		t.Errorf("users nicht in tables.json gefunden")
	}
}

func TestHandleCreateTable_EmptyDBName(t *testing.T) {
	creator := &table.TableCreator{Logger: &logger.NoopLogger{}}
	meta := table.TableMeta{TableName: "users", Fields: []table.FieldMeta{{Name: "id", Type: "string"}}}
	resp := creator.HandleCreateTable("", meta)
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_NAME_MISSING" {
		t.Errorf("Fehlerfall leerer DB-Name nicht erkannt: %+v", resp)
	}
}

func TestHandleCreateTable_InvalidTableName(t *testing.T) {
	creator := &table.TableCreator{Logger: &logger.NoopLogger{}}
	meta := table.TableMeta{TableName: "!invalid", Fields: []table.FieldMeta{{Name: "id", Type: "string"}}}
	resp := creator.HandleCreateTable("testdb", meta)
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_TABLE_INVALID_NAME" {
		t.Errorf("Fehlerfall ungültiger Tabellenname nicht erkannt: %+v", resp)
	}
}

func TestHandleCreateTable_ReservedTableName(t *testing.T) {
	creator := &table.TableCreator{Logger: &logger.NoopLogger{}}
	reserved := []string{"meta", "entries", "indexes"}
	for _, name := range reserved {
		meta := table.TableMeta{TableName: name, Fields: []table.FieldMeta{{Name: "id", Type: "string"}}}
		resp := creator.HandleCreateTable("testdb", meta)
		if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_TABLE_RESERVED_NAME" {
			t.Errorf("Fehlerfall reservierter Tabellenname '%s' nicht erkannt: %+v", name, resp)
		}
	}
}

func TestHandleCreateTable_NoFields(t *testing.T) {
	creator := &table.TableCreator{Logger: &logger.NoopLogger{}}
	meta := table.TableMeta{TableName: "users", Fields: []table.FieldMeta{}}
	resp := creator.HandleCreateTable("testdb", meta)
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_FIELDS_MISSING" {
		t.Errorf("Fehlerfall keine Felder nicht erkannt: %+v", resp)
	}
}

func TestHandleCreateTable_TableExists(t *testing.T) {
	testDir := setupCreateTestDir(t)
	defer teardownCreateTestDir(testDir)
	dbDir := filepath.Join(testDir, "testdb")
	tableDir := filepath.Join(dbDir, "users")
	os.MkdirAll(tableDir, 0755)
	creator := &table.TableCreator{DataDir: testDir, Logger: &logger.NoopLogger{}}
	meta := table.TableMeta{TableName: "users", Fields: []table.FieldMeta{{Name: "id", Type: "string"}}}
	resp := creator.HandleCreateTable("testdb", meta)
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_TABLE_EXISTS" {
		t.Errorf("Fehlerfall Tabelle existiert nicht erkannt: %+v", resp)
	}
}
