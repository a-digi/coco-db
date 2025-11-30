package query_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"strings"
	"github.com/a-digi/coco-db/src/query"
	"github.com/a-digi/coco-db/src/logger"
	"os"
	"path/filepath"
	fields "github.com/a-digi/coco-db/src/table/fields"
)

func TestQueryHandler_GraphQLLikeQuery(t *testing.T) {
	t.Error("Test gestartet")
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	os.RemoveAll("./testdata")
	setupTestTable("./testdata", "testdb", "users", t)
	// Setup orders-Tabelle
	meta := &fields.TableMeta{
		TableName: "orders",
		Fields: []fields.FieldMeta{
			{Name: "id", Type: "string"},
			{Name: "user_id", Type: "string"},
			{Name: "amount", Type: "int"},
			{Name: "status", Type: "string"},
		},
		Indexes: []fields.IndexMeta{
			{Name: "id_idx", Fields: []string{"id"}, Type: "primary", Unique: true},
		},
	}
	metaPath := filepath.Join("./testdata", "testdb", "orders", "meta.json")
	os.MkdirAll(filepath.Dir(metaPath), 0755)
	metaBytes, _ := json.Marshal(meta)
	os.WriteFile(metaPath, metaBytes, 0644)
	entry := map[string]interface{}{"id": "1", "user_id": "1", "amount": 99, "status": "open"}
	entriesDir := filepath.Join("./testdata", "testdb", "orders", "entries")
	os.MkdirAll(entriesDir, 0755)
	entryBytes, _ := json.Marshal(entry)
	os.WriteFile(filepath.Join(entriesDir, "1.json"), entryBytes, 0644)
	idxDir := filepath.Join("./testdata", "testdb", "orders", "indexes")
	os.MkdirAll(idxDir, 0755)
	idxPath := filepath.Join(idxDir, "index_id_idx.json")
	idxObj := map[string][]string{"1": {"1"}}
	idxBytes, _ := json.Marshal(idxObj)
	os.WriteFile(idxPath, idxBytes, 0644)

	graphqlQuery := `query { users( filter: { created_at: { gte: "2025-11-01T00:00:00Z", lte: "2025-11-30T23:59:59Z" }, email: { like: "*@gmail.com" }, age: { gte: 18 } }, limit: 10, offset: 0, sort: ["created_at", "age"], join: [ { table: "orders", on: { user_id: "id" }, filter: { status: "open" }, fields: ["id", "amount", "status"] } ], fields: ["id", "name", "email", "created_at", "age", "orders"] ) { id name email created_at age orders { id amount status } } }`
	r := httptest.NewRequest("POST", "/api/databases/testdb/search", strings.NewReader(graphqlQuery))
	resp := h.QueryHandler("testdb", r)
	if !resp.Success || resp.HttpCode != 200 {
		if resp.Error != nil {
			panic(resp.Error)
		}
		panic("Unbekannter Fehler im GraphQL-Test")
	}
	// Weitere Assertions je nach gewünschtem Ergebnis
}

func setupTestTable(dataDir, dbName, tableName string, t *testing.T) {
	meta := &fields.TableMeta{
		TableName: tableName,
		Fields: []fields.FieldMeta{
			{Name: "id", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
			{Name: "created_at", Type: "string"},
			{Name: "age", Type: "int"},
		},
		Indexes: []fields.IndexMeta{
			{Name: "id_idx", Fields: []string{"id"}, Type: "primary", Unique: true},
		},
	}
	metaPath := filepath.Join(dataDir, dbName, tableName, "meta.json")
	_ = os.MkdirAll(filepath.Dir(metaPath), 0755)
	metaBytes, _ := json.Marshal(meta)
	_ = os.WriteFile(metaPath, metaBytes, 0644)
	// Ein Eintrag, der exakt zu den Filterbedingungen passt
	entry := map[string]interface{}{
		"id": "1",
		"name": "Test User",
		"email": "testuser@gmail.com",
		"created_at": "2025-11-15T12:00:00Z",
		"age": 25,
	}
	entriesDir := filepath.Join(dataDir, dbName, tableName, "entries")
	_ = os.MkdirAll(entriesDir, 0755)
	entryBytes, _ := json.Marshal(entry)
	_ = os.WriteFile(filepath.Join(entriesDir, "1.json"), entryBytes, 0644)
	// Indexdatei als Objekt anlegen (id → [id])
	idxDir := filepath.Join(dataDir, dbName, tableName, "indexes")
	_ = os.MkdirAll(idxDir, 0755)
	idxPath := filepath.Join(idxDir, "index_id_idx.json")
	idxObj := map[string][]string{"1": {"1"}}
	idxBytes, _ := json.Marshal(idxObj)
	_ = os.WriteFile(idxPath, idxBytes, 0644)
}
