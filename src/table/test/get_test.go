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
	tables := []fields.TableMeta{{TableName: tableName, Fields: []fields.FieldMeta{{Name: "id", Type: "string"}}}}
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	f, err := os.Create(tablesJsonPath)
	if err != nil {
		t.Fatalf("Fehler beim Anlegen von tables.json: %v", err)
	}
	_ = json.NewEncoder(f).Encode(tables)
	f.Close()

	tlm := &tbl.ListTablesHandler{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := tlm.HandleListTables(dbName)
	if resp == nil || !resp.Success || resp.HttpCode != 200 {
		t.Fatalf("Erwartet: Success true und HttpCode 200, erhalten: %+v", resp)
	}
	tableList, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Antwortformat unerwartet: %+v", resp.Data)
	}
	tablesIface, ok := tableList["tables"]
	if !ok {
		t.Fatalf("Key 'tables' fehlt in Antwort: %+v", tableList)
	}
	tablesArr, ok := tablesIface.([]string)
	if !ok {
		// Versuche []interface{} zu casten und in []string umzuwandeln (Fallback für manche Go-Versionen)
		if arr, ok2 := tablesIface.([]interface{}); ok2 {
			tablesArr = make([]string, len(arr))
			for i, v := range arr {
				if s, ok3 := v.(string); ok3 {
					tablesArr[i] = s
				}
			}
		} else {
			t.Fatalf("'tables' ist kein []string oder []interface{}: %+v", tablesIface)
		}
	}
	found := false
	for _, t := range tablesArr {
		if t == tableName {
			found = true
		}
	}
	if !found {
		t.Errorf("Tabelle %s nicht in der Liste enthalten: %+v", tableName, tablesArr)
	}
}

func TestHandleGetTable_MissingParams(t *testing.T) {
	testDir := setupGetTestDir(t)
	defer teardownGetTestDir(testDir)
	tlm := &tbl.ListTablesHandler{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := tlm.HandleListTables("")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_NAME_MISSING" {
		t.Errorf("Fehlende Parameter nicht korrekt erkannt: %+v", resp)
	}
}
