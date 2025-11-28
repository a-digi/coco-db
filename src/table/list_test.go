package table

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
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

func TestHandleListTables_Empty(t *testing.T) {
	testDir := setupListTestDir(t)
	defer teardownListTestDir(testDir)
	os.MkdirAll(filepath.Join(testDir, "testdb"), 0755)

	req := httptest.NewRequest(http.MethodGet, "/api/databases/testdb/tables", nil)
	req.Header.Set("X-DB-Name", "testdb")
	w := httptest.NewRecorder()

	h := &ListTablesHandler{
		DataDir:        testDir,
		Logger:         &logger.NoopLogger{},
	}

	h.HandleListTables(w, "testdb")
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status != 200: %d", resp.StatusCode)
	}
	var apiResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&apiResp)
	tables := apiResp["data"].(map[string]interface{})["tables"].([]interface{})
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

	req := httptest.NewRequest(http.MethodGet, "/api/databases/testdb/tables", nil)
	req.Header.Set("X-DB-Name", "testdb")
	w := httptest.NewRecorder()

	h := &ListTablesHandler{
		DataDir:        testDir,
		Logger:         &logger.NoopLogger{},
	}

	h.HandleListTables(w, "testdb")
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status != 200: %d", resp.StatusCode)
	}
	var apiResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&apiResp)
	tables := apiResp["data"].(map[string]interface{})["tables"].([]interface{})
	if len(tables) != 2 {
		t.Errorf("Erwartet: 2 Tabellen, erhalten: %d", len(tables))
	}
}
