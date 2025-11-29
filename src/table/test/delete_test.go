package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/table"
)

func setupDeleteTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_table_delete_test_")
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
	tables := []table.TableMeta{{TableName: tableName, Fields: []table.FieldMeta{{Name: "id", Type: "string"}}}}
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	f, err := os.Create(tablesJsonPath)
	if err != nil {
		t.Fatalf("Fehler beim Anlegen von tables.json: %v", err)
	}
	_ = json.NewEncoder(f).Encode(tables)
	f.Close()

	td := &table.TableDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("DELETE", "/?db=testdb&table=users", nil)
	w := httptest.NewRecorder()
	td.HandleDeleteTable(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Erwartet: Status 200, erhalten: %d", resp.StatusCode)
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
	var tablesMeta []table.TableMeta
	if err := json.Unmarshal(content, &tablesMeta); err != nil {
		t.Fatalf("tables.json nicht parsebar: %v", err)
	}
	for _, entry := range tablesMeta {
		if entry.TableName == tableName {
			t.Errorf("Tabelle wurde nicht aus tables.json entfernt")
		}
	}
}

func TestHandleDeleteTable_NotFound(t *testing.T) {
	testDir := setupDeleteTestDir(t)
	defer teardownDeleteTestDir(testDir)
	td := &table.TableDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("DELETE", "/?db=testdb&table=notfound", nil)
	w := httptest.NewRecorder()
	td.HandleDeleteTable(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Erwartet: Status 404, erhalten: %d", resp.StatusCode)
	}
}

func TestHandleDeleteTable_MissingParams(t *testing.T) {
	testDir := setupDeleteTestDir(t)
	defer teardownDeleteTestDir(testDir)
	td := &table.TableDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("DELETE", "/", nil)
	w := httptest.NewRecorder()
	td.HandleDeleteTable(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Erwartet: Status 400, erhalten: %d", resp.StatusCode)
	}
}

