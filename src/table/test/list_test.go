package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/table"
	"github.com/a-digi/coco-db/src/response"
)

func setupListTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_list_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownListTestDir(dir string) {
	os.RemoveAll(dir)
}

func createTableDir(t *testing.T, dbDir, tableName string) {
	tableDir := filepath.Join(dbDir, tableName)
	os.MkdirAll(tableDir, 0755)
	metaPath := filepath.Join(tableDir, "meta.json")
	os.WriteFile(metaPath, []byte(`{"tableName":"`+tableName+`","fields":[]}`), 0644)
}

func getTablesFromAPIResponse(t *testing.T, resp *response.APIResponse) []string {
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Data fehlt oder hat falschen Typ")
	}
	tablesIface, ok := data["tables"]
	if !ok {
		t.Fatalf("Key 'tables' fehlt in Data")
	}
	tables := []string{}
	switch v := tablesIface.(type) {
	case []interface{}:
		for _, entry := range v {
			if s, ok := entry.(string); ok {
				tables = append(tables, s)
			}
		}
	case []string:
		tables = v
	default:
		t.Fatalf("Tabellen-Array hat unerwarteten Typ: %T", v)
	}
	return tables
}

func TestHandleListTables_Empty(t *testing.T) {
	testDir := setupListTestDir(t)
	defer teardownListTestDir(testDir)
	os.MkdirAll(filepath.Join(testDir, "testdb"), 0755)

	h := &table.ListTablesHandler{
		DataDir: testDir,
		Logger:  &logger.NoopLogger{},
	}

	resp := h.HandleListTables("testdb")
	if resp == nil {
		t.Fatalf("APIResponse ist nil")
	}
	if !resp.Success {
		t.Fatalf("Erwartet: Success true, erhalten: false, Fehler: %v", resp.Error)
	}
	tables := getTablesFromAPIResponse(t, resp)
	if len(tables) != 0 {
		t.Errorf("Erwartet: 0 Tabellen, erhalten: %d", len(tables))
	}
}

func TestHandleListTables_WithTables(t *testing.T) {
	testDir := setupListTestDir(t)
	defer teardownListTestDir(testDir)
	dbDir := filepath.Join(testDir, "testdb")
	os.MkdirAll(dbDir, 0755)
	createTableDir(t, dbDir, "users")
	createTableDir(t, dbDir, "orders")

	h := &table.ListTablesHandler{
		DataDir: testDir,
		Logger:  &logger.NoopLogger{},
	}

	resp := h.HandleListTables("testdb")
	if resp == nil {
		t.Fatalf("APIResponse ist nil")
	}
	if !resp.Success {
		t.Fatalf("Erwartet: Success true, erhalten: false, Fehler: %v", resp.Error)
	}
	tables := getTablesFromAPIResponse(t, resp)
	if len(tables) != 2 {
		t.Errorf("Erwartet: 2 Tabellen, erhalten: %d", len(tables))
	}
}

func TestHandleListTables_ErrorCases(t *testing.T) {
	testDir := setupListTestDir(t)
	defer teardownListTestDir(testDir)

	h := &table.ListTablesHandler{
		DataDir: testDir,
		Logger:  &logger.NoopLogger{},
	}

	// Fehler: Datenbankname fehlt
	resp := h.HandleListTables("")
	if resp == nil || resp.Success || resp.Error == nil {
		t.Fatalf("Fehlerfall: Kein Datenbankname nicht korrekt erkannt")
	}
	if resp.Error.Code != "ERR_DB_NAME_MISSING" {
		t.Errorf("Falscher Fehlercode: %v", resp.Error.Code)
	}
	if resp.HttpCode != 400 {
		t.Errorf("Falscher HttpCode: %v", resp.HttpCode)
	}
	if resp.Error.Message == "" {
		t.Errorf("Fehlermeldung fehlt")
	}

	// Fehler: Datenbank nicht gefunden
	resp = h.HandleListTables("nichtvorhanden")
	if resp == nil || resp.Success || resp.Error == nil {
		t.Fatalf("Fehlerfall: Nicht vorhandene DB nicht korrekt erkannt")
	}
	if resp.Error.Code != "ERR_DB_NOT_FOUND" {
		t.Errorf("Falscher Fehlercode: %v", resp.Error.Code)
	}
	if resp.HttpCode != 404 {
		t.Errorf("Falscher HttpCode: %v", resp.HttpCode)
	}
	if resp.Error.Message == "" {
		t.Errorf("Fehlermeldung fehlt")
	}
}
