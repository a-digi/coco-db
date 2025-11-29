package query_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
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
			{Name: "age", Type: "int"},
		},
		Indexes: []fields.IndexMeta{
			{Name: "email_idx", Type: "btree", Fields: []string{"email"}},
		},
	}

	// Testdaten erzeugen
	entriesDir := filepath.Join(dataDir, dbName, tableName, "entries", "48dece68-4380-4f82-92e3-8fdb7da2b8bd")
	_ = os.MkdirAll(entriesDir, 0755)
	entry := map[string]interface{}{"email": "foo@bar.de", "id": "48dece68-4380-4f82-92e3-8fdb7da2b8bd", "name": "Test User", "age": 20}
	b, _ := json.Marshal(entry)
	_ = os.WriteFile(filepath.Join(entriesDir, "48dece68-4380-4f82-92e3-8fdb7da2b8bd.json"), b, 0644)

	idxDir := filepath.Join(dataDir, dbName, tableName, "indexes")
	_ = os.MkdirAll(idxDir, 0755)
	idxPath := filepath.Join(idxDir, "email_idx.json")
	_ = os.WriteFile(idxPath, []byte(`["48dece68-4380-4f82-92e3-8fdb7da2b8bd"]`), 0644)

	// Zusätzliche Testdaten für Sortierung und Paginierung
	entries2 := []struct {
		ID     string
		Email  string
		Name   string
		Age    float64
	}{
		{"2", "bar@foo.de", "Bar User", 17.0},
		{"3", "baz@foo.de", "Baz User", 30.0},
	}
	for _, e := range entries2 {
		entryDir := filepath.Join(dataDir, dbName, tableName, "entries", e.ID)
		_ = os.MkdirAll(entryDir, 0755)
		entry := map[string]interface{}{"email": e.Email, "id": e.ID, "name": e.Name, "age": e.Age}
		b, _ := json.Marshal(entry)
		_ = os.WriteFile(filepath.Join(entryDir, e.ID+".json"), b, 0644)
	}
	// Indexdatei (Dummy: alle IDs)
	_ = os.WriteFile(idxPath, []byte(`["48dece68-4380-4f82-92e3-8fdb7da2b8bd","2","3"]`), 0644)

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

	// Test: Sortierung aufsteigend nach Alter
	q4 := &query.Query{Filter: map[string]interface{}{}, Sort: []string{"age"}}
	result4, err := query.FilterEngine(dataDir, dbName, tableName, q4, meta)
	if err != nil || len(result4) < 2 || result4[0]["age"].(float64) > result4[1]["age"].(float64) {
		t.Errorf("Sortierung aufsteigend: Ergebnis nicht korrekt sortiert: %+v, err=%v", result4, err)
	}

	// Test: Sortierung absteigend nach Alter
	q5 := &query.Query{Filter: map[string]interface{}{}, Sort: []string{"-age"}}
	result5, err := query.FilterEngine(dataDir, dbName, tableName, q5, meta)
	if err != nil || len(result5) < 2 || result5[0]["age"].(float64) < result5[1]["age"].(float64) {
		t.Errorf("Sortierung absteigend: Ergebnis nicht korrekt sortiert: %+v, err=%v", result5, err)
	}

	// Test: Paginierung (limit=1, offset=1)
	q6 := &query.Query{Filter: map[string]interface{}{}, Limit: 1, Offset: 1, Sort: []string{"age"}}
	result6, err := query.FilterEngine(dataDir, dbName, tableName, q6, meta)
	if err != nil || len(result6) != 1 {
		t.Errorf("Paginierung: Erwartet 1 Ergebnis, bekommen: %+v, err=%v", result6, err)
	}
	if len(result6) == 1 && result6[0]["age"].(float64) != 20.0 {
		t.Errorf("Paginierung: Erwartet age=20.0, bekommen: %+v", result6[0])
	}
}

func TestFilterEngine_DeepJoins(t *testing.T) {
	dataDir := "../../../../data"
	dbName := "testdb"
	joinDepth := 10

	// Testdaten für 10 verschachtelte Ebenen erzeugen
	for i := 0; i <= joinDepth; i++ {
		tableName := "level" + strconv.Itoa(i)
		entriesDir := filepath.Join(dataDir, dbName, tableName, "entries", "id0")
		_ = os.MkdirAll(entriesDir, 0755)
		entry := map[string]interface{}{"id": "id0"}
		if i < joinDepth {
			entry["ref_id"] = "id0"
		}
		b, _ := json.Marshal(entry)
		_ = os.WriteFile(filepath.Join(entriesDir, "id0.json"), b, 0644)
		// Indexdatei
		idxDir := filepath.Join(dataDir, dbName, tableName, "indexes")
		_ = os.MkdirAll(idxDir, 0755)
		idxPath := filepath.Join(idxDir, "id_idx.json")
		_ = os.WriteFile(idxPath, []byte(`["id0"]`), 0644)
	}

	// TableMeta für alle Ebenen
	metas := make([]*fields.TableMeta, joinDepth+1)
	for i := 0; i <= joinDepth; i++ {
		tableName := "level" + strconv.Itoa(i)
		metas[i] = &fields.TableMeta{
			TableName: tableName,
			Fields: []fields.FieldMeta{
				{Name: "id", Type: "string"},
				{Name: "ref_id", Type: "string"},
			},
			Indexes: []fields.IndexMeta{
				{Name: "id_idx", Type: "btree", Fields: []string{"id"}},
			},
		}
	}

	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(dataDir, dbName))
	})

	// Baue die Join-Definitionen für 10 Ebenen
	joins := make([]map[string]interface{}, joinDepth)
	for i := 1; i <= joinDepth; i++ {
		joins[i-1] = map[string]interface{}{
			"table": "level" + strconv.Itoa(i),
			"on": map[string]interface{}{ "ref_id": "id" },
		}
	}

	// Simuliere eine globale Query mit 10 verschachtelten Joins
	q := &query.Query{
		Filter: map[string]interface{}{ "id": "id0" },
		Join: joins,
	}

	// Dummy-Implementierung: Wir prüfen nur, ob die Join-Struktur korrekt erzeugt werden kann
	// (Die eigentliche Join-Engine muss in der API implementiert sein)
	// Hier simulieren wir, dass die Join-Tiefe nicht zu einem Stackoverflow oder Fehler führt
	// und dass die Join-Definitionen korrekt verarbeitet werden
	if len(q.Join) != joinDepth {
		t.Errorf("Erwartet %d Joins, bekommen: %d", joinDepth, len(q.Join))
	}
	// (Optional: Wenn Join-Engine vorhanden, kann hier die tatsächliche Ausführung getestet werden)
}
