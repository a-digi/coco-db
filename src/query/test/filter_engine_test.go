package query_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/query"
	"github.com/a-digi/coco-db/src/table/fields"
)

func TestFilterEngine_IndexAndNonIndex(t *testing.T) {
	dataDir := "../../../../data"
	dbName := "testdb"
	tableName := "users"
	meta := &fields.TableMeta{
		TableName: tableName,
		Fields: []fields.FieldMeta{
			{Name: "email", Type: "string"},
			{Name: "name", Type: "string"},
		},
		Indexes: []fields.IndexMeta{
			{Name: "email_idx", Type: "btree", Fields: []string{"email"}},
		},
	}

	// Testdaten erzeugen
	entriesDir := filepath.Join(dataDir, dbName, tableName, "entries", "48dece68-4380-4f82-92e3-8fdb7da2b8bd")
	_ = os.MkdirAll(entriesDir, 0755)
	entry := map[string]interface{}{"email": "foo@bar.de", "id": "48dece68-4380-4f82-92e3-8fdb7da2b8bd", "name": "Test User"}
	b, _ := json.Marshal(entry)
	_ = os.WriteFile(filepath.Join(entriesDir, "48dece68-4380-4f82-92e3-8fdb7da2b8bd.json"), b, 0644)

	idxDir := filepath.Join(dataDir, dbName, tableName, "indexes")
	_ = os.MkdirAll(idxDir, 0755)
	idxPath := filepath.Join(idxDir, "email_idx.json")
	_ = os.WriteFile(idxPath, []byte(`["48dece68-4380-4f82-92e3-8fdb7da2b8bd"]`), 0644)

	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(dataDir, dbName))
	})

	// Test: Filter auf Indexfeld (email)
	q := &query.Query{Filter: map[string]interface{}{"email": "foo@bar.de"}}
	result, err := query.FilterEngine(dataDir, dbName, tableName, q, meta)
	if err != nil {
		t.Fatalf("Fehler bei FilterEngine (Index): %v", err)
	}
	if len(result) != 1 || result[0]["email"] != "foo@bar.de" {
		t.Errorf("Index-Filter: Erwartet 1 Ergebnis mit email=foo@bar.de, bekommen: %+v", result)
	}

	// Test: Filter auf nicht indiziertes Feld (name)
	q2 := &query.Query{Filter: map[string]interface{}{"name": "Test User"}}
	result2, err := query.FilterEngine(dataDir, dbName, tableName, q2, meta)
	if err != nil {
		t.Fatalf("Fehler bei FilterEngine (Non-Index): %v", err)
	}
	if len(result2) != 1 || result2[0]["name"] != "Test User" {
		t.Errorf("Non-Index-Filter: Erwartet 1 Ergebnis mit name=Test User, bekommen: %+v", result2)
	}

	// Test: AND-Logik (beide Felder)
	q3 := &query.Query{Filter: map[string]interface{}{"email": "foo@bar.de", "name": "Test User"}}
	result3, err := query.FilterEngine(dataDir, dbName, tableName, q3, meta)
	if err != nil {
		t.Fatalf("Fehler bei FilterEngine (AND): %v", err)
	}
	if len(result3) != 1 || result3[0]["email"] != "foo@bar.de" || result3[0]["name"] != "Test User" {
		t.Errorf("AND-Logik: Erwartet 1 Ergebnis mit email=foo@bar.de und name=Test User, bekommen: %+v", result3)
	}
}
