package query_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/query"
	"github.com/a-digi/coco-db/src/table/fields"
)

func setupTestData(t *testing.T) (dataDir, dbName, tableName string, meta *fields.TableMeta, cleanup func()) {
	dataDir = "../../../../data"
	dbName = "testdb"
	tableName = "users"
	meta = &fields.TableMeta{
		TableName: tableName,
		Fields: []fields.FieldMeta{
			{Name: "email", Type: "string"},
			{Name: "age", Type: "int"},
			{Name: "isActive", Type: "bool"},
		},
		Indexes: []fields.IndexMeta{
			{Name: "email_idx", Type: "btree", Fields: []string{"email"}},
		},
	}
	// Testdaten erzeugen
	entries := []struct {
		ID     string
		Email  string
		Age    float64
		Active bool
	}{
		{"1", "foo@bar.de", 20.0, true},
		{"2", "bar@foo.de", 17.0, false},
		{"3", "baz@foo.de", 30.0, true},
	}
	for _, e := range entries {
		entryDir := filepath.Join(dataDir, dbName, tableName, "entries", e.ID)
		_ = os.MkdirAll(entryDir, 0755)
		entry := map[string]interface{}{"email": e.Email, "age": e.Age, "isActive": e.Active, "id": e.ID}
		b, _ := json.Marshal(entry)
		_ = os.WriteFile(filepath.Join(entryDir, e.ID+".json"), b, 0644)
	}
	// Indexdatei (Dummy: alle IDs)
	idxDir := filepath.Join(dataDir, dbName, tableName, "indexes")
	_ = os.MkdirAll(idxDir, 0755)
	idxPath := filepath.Join(idxDir, "email_idx.json")
	_ = os.WriteFile(idxPath, []byte(`["1","2","3"]`), 0644)

	cleanup = func() {
		os.RemoveAll(filepath.Join(dataDir, dbName))
	}
	return
}

func TestFilterEngine_Basic(t *testing.T) {
	dataDir, dbName, tableName, meta, cleanup := setupTestData(t)
	defer cleanup()

	t.Run("eq filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"email": "foo@bar.de"}}
		res, err := query.FilterEngine(dataDir, dbName, tableName, q, meta)
		if err != nil || len(res) != 1 || res[0]["email"] != "foo@bar.de" {
			t.Errorf("expected 1 result, got %+v, err=%v", res, err)
		}
	})

	t.Run("gt filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"age": map[string]interface{}{"gt": 18}}}
		res, err := query.FilterEngine(dataDir, dbName, tableName, q, meta)
		if err != nil || len(res) != 2 {
			t.Errorf("expected 2 results, got %+v, err=%v", res, err)
		}
	})

	t.Run("and logic", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"isActive": true, "age": map[string]interface{}{"gte": 18}}}
		res, err := query.FilterEngine(dataDir, dbName, tableName, q, meta)
		if err != nil || len(res) != 2 {
			t.Errorf("expected 2 results, got %+v, err=%v", res, err)
		}
	})

	t.Run("like filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"email": map[string]interface{}{"like": "foo*@bar.de"}}}
		res, err := query.FilterEngine(dataDir, dbName, tableName, q, meta)
		if err != nil || len(res) != 1 || res[0]["email"] != "foo@bar.de" {
			t.Errorf("expected 1 result for LIKE, got %+v, err=%v", res, err)
		}
	})

	t.Run("partial filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"email": map[string]interface{}{"partial": "foo"}}}
		res, err := query.FilterEngine(dataDir, dbName, tableName, q, meta)
		if err != nil || len(res) != 2 {
			t.Errorf("expected 2 results for PARTIAL, got %+v, err=%v", res, err)
		}
	})

	t.Run("fulltext filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"email": map[string]interface{}{"fulltext": "foo bar"}}}
		res, err := query.FilterEngine(dataDir, dbName, tableName, q, meta)
		if err != nil || len(res) != 1 || res[0]["email"] != "foo@bar.de" {
			t.Errorf("expected 1 result for FULLTEXT, got %+v, err=%v", res, err)
		}
	})
}
