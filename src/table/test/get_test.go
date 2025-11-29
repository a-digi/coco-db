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
	r := httptest.NewRequest("GET", "/?db=testdb&table=users", nil)
	w := httptest.NewRecorder()
	tlm.HandleGetTable(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Erwartet: Status 200, erhalten: %d", resp.StatusCode)
	}
	var apiResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Antwort nicht parsebar: %v", err)
	}
	if success, ok := apiResp["success"].(bool); !ok || !success {
		t.Errorf("Erwartet: success true, erhalten: %+v", apiResp)
	}
	if data, ok := apiResp["data"].(map[string]interface{}); !ok || data["tableName"] != tableName {
		t.Errorf("Tabellenname nicht korrekt: %+v", apiResp)
	}
}

func TestHandleGetTable_NotFound(t *testing.T) {
	testDir := setupGetTestDir(t)
	defer teardownGetTestDir(testDir)
	tlm := &table.TableListMeta{DataDir: testDir, Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("GET", "/?db=testdb&table=notfound", nil)
	w := httptest.NewRecorder()
	tlm.HandleGetTable(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Erwartet: Status 404, erhalten: %d", resp.StatusCode)
	}
	var apiResp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&apiResp)
	if success, ok := apiResp["success"].(bool); ok && success {
		t.Errorf("Erwartet: success false, erhalten: %+v", apiResp)
	}
}

func TestHandleGetTable_MissingParams(t *testing.T) {
	testDir := setupGetTestDir(t)
	defer teardownGetTestDir(testDir)
	tlm := &table.TableListMeta{DataDir: testDir, Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	tlm.HandleGetTable(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Erwartet: Status 400, erhalten: %d", resp.StatusCode)
	}
	var apiResp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&apiResp)
	if success, ok := apiResp["success"].(bool); ok && success {
		t.Errorf("Erwartet: success false, erhalten: %+v", apiResp)
	}
}
