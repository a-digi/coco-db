package table_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/table"
	"github.com/a-digi/coco-db/src/table/fields"
)

func setupEditTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_edit_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownEditTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestHandleEditTable_Success(t *testing.T) {
	testDir := setupEditTestDir(t)
	defer teardownEditTestDir(testDir)
	dbName := "testdb"
	tableName := "users"
	dbDir := filepath.Join(testDir, dbName)
	tableDir := filepath.Join(dbDir, tableName)
	os.MkdirAll(tableDir, 0755)
	// Lege initiale meta.json an
	initMeta := fields.TableMeta{
		TableName: tableName,
		Fields:    []fields.FieldMeta{{Name: "id", Type: "string"}},
	}
	metaPath := filepath.Join(tableDir, "meta.json")
	f, err := os.Create(metaPath)
	if err != nil {
		t.Fatalf("Fehler beim Anlegen von meta.json: %v", err)
	}
	_ = json.NewEncoder(f).Encode(initMeta)
	f.Close()

	tu := &table.TableUpdate{DataDir: testDir, Logger: &logger.NoopLogger{}}
	newMeta := fields.TableMeta{
		TableName: tableName,
		Fields:    []fields.FieldMeta{{Name: "id", Type: "string"}, {Name: "email", Type: "string"}},
	}
	resp := tu.HandleEditTable(dbName, tableName, newMeta)
	if resp == nil || !resp.Success || resp.HttpCode != 200 {
		t.Fatalf("Erwartet: Success true und HttpCode 200, erhalten: %+v", resp)
	}
	// Prüfe, ob meta.json aktualisiert wurde
	content, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("meta.json nicht lesbar: %v", err)
	}
	var gotMeta fields.TableMeta
	if err := json.Unmarshal(content, &gotMeta); err != nil {
		t.Fatalf("meta.json nicht parsebar: %v", err)
	}
	if len(gotMeta.Fields) != 2 {
		t.Errorf("meta.json nicht aktualisiert: %+v", gotMeta)
	}
	// Prüfe, ob tables.json aktualisiert wurde
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	tablesContent, err := os.ReadFile(tablesJsonPath)
	if err != nil {
		t.Fatalf("tables.json nicht lesbar: %v", err)
	}
	var tablesMeta []fields.TableMeta
	if err := json.Unmarshal(tablesContent, &tablesMeta); err != nil {
		t.Fatalf("tables.json nicht parsebar: %v", err)
	}
	found := false
	for _, entry := range tablesMeta {
		if entry.TableName == tableName && len(entry.Fields) == 2 {
			found = true
		}
	}
	if !found {
		t.Errorf("Tabelle nicht korrekt in tables.json aktualisiert")
	}
}

func TestHandleEditTable_TableNotFound(t *testing.T) {
	testDir := setupEditTestDir(t)
	defer teardownEditTestDir(testDir)
	tu := &table.TableUpdate{DataDir: testDir, Logger: &logger.NoopLogger{}}
	meta := fields.TableMeta{
		TableName: "users",
		Fields:    []fields.FieldMeta{{Name: "id", Type: "string"}},
	}
	resp := tu.HandleEditTable("testdb", "users", meta)
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_TABLE_NOT_FOUND" {
		t.Errorf("Nicht vorhandene Tabelle nicht korrekt erkannt: %+v", resp)
	}
}

func TestHandleEditTable_InvalidFields(t *testing.T) {
	testDir := setupEditTestDir(t)
	defer teardownEditTestDir(testDir)
	dbName := "testdb"
	tableName := "users"
	dbDir := filepath.Join(testDir, dbName)
	tableDir := filepath.Join(dbDir, tableName)
	os.MkdirAll(tableDir, 0755)
	initMeta := fields.TableMeta{
		TableName: tableName,
		Fields:    []fields.FieldMeta{{Name: "id", Type: "string"}},
	}
	metaPath := filepath.Join(tableDir, "meta.json")
	f, err := os.Create(metaPath)
	if err != nil {
		t.Fatalf("Fehler beim Anlegen von meta.json: %v", err)
	}
	_ = json.NewEncoder(f).Encode(initMeta)
	f.Close()

	tu := &table.TableUpdate{DataDir: testDir, Logger: &logger.NoopLogger{}}
	invalidMeta := fields.TableMeta{
		TableName: tableName,
		Fields:    []fields.FieldMeta{}, // keine Felder
	}
	resp := tu.HandleEditTable(dbName, tableName, invalidMeta)
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_FIELD_INVALID" {
		t.Errorf("Ungültige Felder nicht korrekt erkannt: %+v", resp)
	}
}
