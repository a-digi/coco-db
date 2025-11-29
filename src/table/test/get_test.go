package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/table"
)

func setupGetTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_get_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownGetTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestHandleGetTable_Success(t *testing.T) {
	testDir := setupGetTestDir(t)
	defer teardownGetTestDir(testDir)
	dbName := "testdb"
	tableName := "users"
	dbDir := filepath.Join(testDir, dbName)
	os.MkdirAll(dbDir, 0755)
	// Lege tables.json mit einer Tabelle an
	tables := []table.TableMeta{{TableName: tableName, Fields: []table.FieldMeta{{Name: "id", Type: "string"}}}}
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	f, err := os.Create(tablesJsonPath)
	if err != nil {
		t.Fatalf("Fehler beim Anlegen von tables.json: %v", err)
	}
	_ = json.NewEncoder(f).Encode(tables)
	f.Close()

	tlm := &table.TableListMeta{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := tlm.HandleGetTable(dbName, tableName)
	if resp == nil || !resp.Success || resp.HttpCode != 200 {
		t.Fatalf("Erwartet: Success true und HttpCode 200, erhalten: %+v", resp)
	}
	meta, ok := resp.Data.(table.TableMeta)
	if !ok || meta.TableName != tableName {
		t.Errorf("Tabellenname nicht korrekt: %+v", resp.Data)
	}
}

func TestHandleGetTable_TableNotFound(t *testing.T) {
	testDir := setupGetTestDir(t)
	defer teardownGetTestDir(testDir)
	tlm := &table.TableListMeta{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := tlm.HandleGetTable("testdb", "notfound")
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

func TestHandleGetTable_MissingParams(t *testing.T) {
	testDir := setupGetTestDir(t)
	defer teardownGetTestDir(testDir)
	tlm := &table.TableListMeta{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := tlm.HandleGetTable("", "")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_PARAM_MISSING" {
		t.Errorf("Fehlende Parameter nicht korrekt erkannt: %+v", resp)
	}
}
