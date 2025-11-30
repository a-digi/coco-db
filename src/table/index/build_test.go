package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	fields "github.com/a-digi/coco-db/src/table/fields"
)

func TestBuildIndexesFromEntries_SingleField(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "indexes"), 0755)
	entriesDir := filepath.Join(dir, "entries")
	os.MkdirAll(entriesDir, 0755)
	// Lege zwei Einträge an
	entry1 := map[string]interface{}{ "id": "id1", "email": "foo@bar.de" }
	entry2 := map[string]interface{}{ "id": "id2", "email": "bar@baz.de" }
	for _, e := range []map[string]interface{}{entry1, entry2} {
		eid := e["id"].(string)
		eDir := filepath.Join(entriesDir, eid)
		os.MkdirAll(eDir, 0755)
		f, _ := os.Create(filepath.Join(eDir, eid+".json"))
		_ = json.NewEncoder(f).Encode(e)
		f.Close()
	}
	meta := fields.TableMeta{
		TableName: "users",
		Fields: []fields.FieldMeta{{Name: "id", Type: "string"}, {Name: "email", Type: "string"}},
		Indexes: []fields.IndexMeta{{Name: "email_idx", Type: "secondary", Fields: []string{"email"}, Unique: false}},
	}
	err := BuildIndexesFromEntries(dir, meta)
	if err != nil {
		t.Fatalf("Fehler beim Indexaufbau: %v", err)
	}
	// Prüfe, ob Indexdatei existiert und beide Einträge enthält
	indexPath := filepath.Join(dir, "indexes", "index_email_idx.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Indexdatei nicht gefunden: %v", err)
	}
	var idxMap map[string][]string
	if err := json.Unmarshal(data, &idxMap); err != nil {
		t.Fatalf("Indexdatei nicht parsebar: %v", err)
	}
	if len(idxMap) != 2 || len(idxMap["foo@bar.de"]) != 1 || len(idxMap["bar@baz.de"]) != 1 {
		t.Errorf("Indexdaten stimmen nicht: %+v", idxMap)
	}
}
