package table

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
)

func setupTestDataDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "cocodb_test_data_")
	if err != nil {
		t.Fatalf("Fehler beim Anlegen des Testverzeichnisses: %v", err)
	}
	return dir
}

func teardownTestDataDir(dir string) {
	os.RemoveAll(dir)
}

func newTestTableCreator(dataDir string) *TableCreator {
	return &TableCreator{
		DataDir: dataDir,
		Logger:  &logger.NoopLogger{},
	}
}

func TestHandleCreateTable_Success(t *testing.T) {
	testDir := setupTestDataDir(t)
	defer teardownTestDataDir(testDir)

	dbname := "testdb"
	tablename := "users"
	meta := TableMeta{
		TableName: tablename,
		Fields:    []FieldMeta{{Name: "id", Type: "string", Required: true}},
	}
	body, _ := json.Marshal(meta)
	req := httptest.NewRequest(http.MethodPost, "/api/databases/testdb/tables", bytes.NewReader(body))
	req.Header.Set("X-DB-Name", dbname)

	tc := newTestTableCreator(testDir)
	w := httptest.NewRecorder()
	tc.HandleCreateTable(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status != 200: %d", resp.StatusCode)
	}
	var apiResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&apiResp)
	if apiResp["success"] != true {
		t.Error("APIResponse success != true")
	}
	// Prüfe, ob meta.json existiert
	metaPath := filepath.Join(testDir, dbname, tablename, "meta.json")
	if _, err := os.Stat(metaPath); err != nil {
		t.Errorf("meta.json wurde nicht angelegt: %v", err)
	}
}

func TestHandleCreateTable_InvalidDBName(t *testing.T) {
	tc := newTestTableCreator("/tmp/irgendwas")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/databases/xxx/tables", strings.NewReader("{}"))
	req.Header.Set("X-DB-Name", "!!invalid!!")
	tc.HandleCreateTable(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status != 400 bei ungültigem DB-Namen: %d", w.Code)
	}
}

func TestHandleCreateTable_InvalidJSON(t *testing.T) {
	tc := newTestTableCreator("/tmp/irgendwas")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/databases/testdb/tables", strings.NewReader("{invalid json}"))
	req.Header.Set("X-DB-Name", "testdb")
	tc.HandleCreateTable(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status != 400 bei ungültigem JSON: %d", w.Code)
	}
}

func TestHandleCreateTable_ReservedTableName(t *testing.T) {
	tc := newTestTableCreator("/tmp/irgendwas")
	w := httptest.NewRecorder()
	meta := TableMeta{TableName: "meta", Fields: []FieldMeta{{Name: "id", Type: "string"}}}
	body, _ := json.Marshal(meta)
	req := httptest.NewRequest(http.MethodPost, "/api/databases/testdb/tables", bytes.NewReader(body))
	req.Header.Set("X-DB-Name", "testdb")
	tc.HandleCreateTable(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status != 400 bei reserviertem Tabellennamen: %d", w.Code)
	}
}

func TestHandleCreateTable_MissingFields(t *testing.T) {
	tc := newTestTableCreator("/tmp/irgendwas")
	w := httptest.NewRecorder()
	meta := TableMeta{TableName: "users"}
	body, _ := json.Marshal(meta)
	req := httptest.NewRequest(http.MethodPost, "/api/databases/testdb/tables", bytes.NewReader(body))
	req.Header.Set("X-DB-Name", "testdb")
	tc.HandleCreateTable(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status != 400 bei fehlenden Feldern: %d", w.Code)
	}
}

func TestHandleCreateTable_FieldWithoutNameOrType(t *testing.T) {
	tc := newTestTableCreator("/tmp/irgendwas")
	w := httptest.NewRecorder()
	meta := TableMeta{TableName: "users", Fields: []FieldMeta{{Name: "", Type: "string"}}}
	body, _ := json.Marshal(meta)
	req := httptest.NewRequest(http.MethodPost, "/api/databases/testdb/tables", bytes.NewReader(body))
	req.Header.Set("X-DB-Name", "testdb")
	tc.HandleCreateTable(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status != 400 bei Feld ohne Namen: %d", w.Code)
	}
}

func TestHandleCreateTable_TableExists(t *testing.T) {
	testDir := setupTestDataDir(t)
	defer teardownTestDataDir(testDir)
	dbname := "testdb"
	tablename := "users"
	os.MkdirAll(filepath.Join(testDir, dbname, tablename), 0755)
	meta := TableMeta{TableName: tablename, Fields: []FieldMeta{{Name: "id", Type: "string"}}}
	body, _ := json.Marshal(meta)
	req := httptest.NewRequest(http.MethodPost, "/api/databases/testdb/tables", bytes.NewReader(body))
	req.Header.Set("X-DB-Name", dbname)
	tc := newTestTableCreator(testDir)
	w := httptest.NewRecorder()
	tc.HandleCreateTable(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Status != 409 bei existierender Tabelle: %d", w.Code)
	}
}
