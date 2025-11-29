package query_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
	"path/filepath"

	"github.com/a-digi/coco-db/src/query"
	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/table/fields"
)

func setupTestTable(dataDir, dbName, tableName string, t *testing.T) {
	meta := &fields.TableMeta{
		TableName: tableName,
		Fields: []fields.FieldMeta{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "age", Type: "int"},
		},
	}
	metaPath := filepath.Join(dataDir, dbName, tableName, "meta.json")
	_ = os.MkdirAll(filepath.Dir(metaPath), 0755)
	metaBytes, _ := json.Marshal(meta)
	_ = os.WriteFile(metaPath, metaBytes, 0644)
	// Ein Eintrag
	entryDir := filepath.Join(dataDir, dbName, tableName, "entries", "1")
	_ = os.MkdirAll(entryDir, 0755)
	entry := map[string]interface{}{"id": "1", "name": "Test", "age": 42}
	entryBytes, _ := json.Marshal(entry)
	_ = os.WriteFile(filepath.Join(entryDir, "1.json"), entryBytes, 0644)
}

func TestTableQueryHandler_Success(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	setupTestTable("./testdata", "testdb", "users", t)
	t.Cleanup(func() { os.RemoveAll("./testdata") })
	queryBody := map[string]interface{}{
		"filter": map[string]interface{}{"id": "1"},
	}
	body, _ := json.Marshal(queryBody)
	r := httptest.NewRequest("POST", "/api/databases/testdb/tables/users/query", bytes.NewReader(body))
	resp := h.TableQueryHandler("testdb", "users", r)
	if !resp.Success || resp.HttpCode != http.StatusOK {
		t.Errorf("Erwartet Success und Status 200, bekommen: %+v", resp)
	}
	data, ok := resp.Data.([]map[string]interface{})
	if !ok || len(data) != 1 || data[0]["id"] != "1" {
		t.Errorf("Erwartet einen Eintrag mit id=1, bekommen: %+v", resp.Data)
	}
}

func TestTableQueryHandler_InvalidJSON(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("POST", "/api/databases/testdb/tables/users/query", bytes.NewReader([]byte("{")))
	resp := h.TableQueryHandler("testdb", "users", r)
	if resp.Success || resp.HttpCode != http.StatusBadRequest {
		t.Errorf("Erwartet Fehler bei ungültigem JSON, bekommen: %+v", resp)
	}
}

func TestTableQueryHandler_NotFound(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	// Kein Setup, Tabelle existiert nicht
	r := httptest.NewRequest("POST", "/api/databases/testdb/tables/users/query", bytes.NewReader([]byte(`{"filter":{"id":"1"}}`)))
	resp := h.TableQueryHandler("testdb", "users", r)
	if resp.Success || resp.HttpCode != http.StatusNotFound {
		t.Errorf("Erwartet Fehler bei fehlender Tabelle, bekommen: %+v", resp)
	}
}

func TestQueryHandler_Success(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	setupTestTable("./testdata", "testdb", "users", t)
	t.Cleanup(func() { os.RemoveAll("./testdata") })
	queryBody := map[string]interface{}{
		"filter": map[string]interface{}{"id": "1"},
	}
	body, _ := json.Marshal(queryBody)
	r := httptest.NewRequest("POST", "/api/query", bytes.NewReader(body))
	resp := h.QueryHandler("testdb", r)
	if !resp.Success || resp.HttpCode != http.StatusOK {
		t.Errorf("Erwartet Success und Status 200, bekommen: %+v", resp)
	}
	data, ok := resp.Data.([]map[string]interface{})
	if !ok || len(data) != 1 || data[0]["id"] != "1" {
		t.Errorf("Erwartet einen Eintrag mit id=1, bekommen: %+v", resp.Data)
	}
}

func TestQueryHandler_NotFound(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("POST", "/api/query", bytes.NewReader([]byte(`{"filter":{"id":"1"}}`)))
	resp := h.QueryHandler("testdb", r)
	if resp.Success || resp.HttpCode != http.StatusNotFound {
		t.Errorf("Erwartet Fehler bei fehlender Tabelle, bekommen: %+v", resp)
	}
}
