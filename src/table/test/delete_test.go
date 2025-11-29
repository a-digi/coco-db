package table_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	tbl "github.com/a-digi/coco-db/src/table"
	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/table/fields"
)

func setupDeleteTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_delete_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownDeleteTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestHandleDeleteTable_Success(t *testing.T) {
	testDir := setupDeleteTestDir(t)
	defer teardownDeleteTestDir(testDir)
	dbName := "testdb"
	tableName := "users"
	dbDir := filepath.Join(testDir, dbName)
	tableDir := filepath.Join(dbDir, tableName)
	os.MkdirAll(tableDir, 0755)
	// Lege tables.json mit der Tabelle an
	tables := []fields.TableMeta{{TableName: tableName, Fields: []fields.FieldMeta{{Name: "id", Type: "string"}}}}
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	f, err := os.Create(tablesJsonPath)
	if err != nil {
		t.Fatalf("Fehler beim Anlegen von tables.json: %v", err)
	}
	_ = json.NewEncoder(f).Encode(tables)
	f.Close()

	td := &tbl.TableDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := td.HandleDeleteTable(dbName, tableName)
	if resp == nil || !resp.Success || resp.HttpCode != 200 {
		t.Fatalf("Erwartet: Success true und HttpCode 200, erhalten: %+v", resp)
	}
	// Prüfe, ob das Tabellenverzeichnis gelöscht wurde
	if _, err := os.Stat(tableDir); !os.IsNotExist(err) {
		t.Errorf("Tabellenverzeichnis wurde nicht gelöscht: %v", err)
	}
	// Prüfe, ob die Tabelle aus tables.json entfernt wurde
	content, err := os.ReadFile(tablesJsonPath)
	if err != nil {
		t.Fatalf("tables.json nicht lesbar: %v", err)
	}
	var tablesMeta []fields.TableMeta
	if err := json.Unmarshal(content, &tablesMeta); err != nil {
		t.Fatalf("tables.json nicht parsebar: %v", err)
	}
	for _, entry := range tablesMeta {
		if entry.TableName == tableName {
			t.Errorf("Tabelle wurde nicht aus tables.json entfernt")
		}
	}
}

func TestHandleDeleteTable_TableNotFound(t *testing.T) {
	testDir := setupDeleteTestDir(t)
	defer teardownDeleteTestDir(testDir)
	td := &tbl.TableDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := td.HandleDeleteTable("testdb", "notfound")
	if resp == nil || resp.Success || resp.Error == nil {
		t.Errorf("Nicht vorhandene Tabelle nicht korrekt erkannt: %+v", resp)
	}
	if resp.Error != nil {
		if resp.Error.Code != "ERR_TABLE_NOT_FOUND" && resp.Error.Code != "ERR_TABLES_NOT_FOUND" {
			t.Errorf("Falscher Fehlercode: %v", resp.Error.Code)
		}
		if resp.HttpCode != 404 {
			t.Errorf("Falscher HttpCode: %v", resp.HttpCode)
		}
		if resp.Error.Message == "" {
			t.Errorf("Fehlermeldung fehlt: %+v", resp.Error)
		}
	}
}

func TestHandleDeleteTable_MissingParams(t *testing.T) {
	testDir := setupDeleteTestDir(t)
	defer teardownDeleteTestDir(testDir)
	td := &tbl.TableDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := td.HandleDeleteTable("", "")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_PARAM_MISSING" {
		t.Errorf("Fehlende Parameter nicht korrekt erkannt: %+v", resp)
	}
}
