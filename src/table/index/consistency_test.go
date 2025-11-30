package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	fields "github.com/a-digi/coco-db/src/table/fields"
)

func TestCheckIndexConsistency_SingleField(t *testing.T) {
	dir := t.TempDir()
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
	// Lege Indexdatei an (nur einen Eintrag absichtlich)
	idxObj := map[string][]string{"foo@bar.de": {"id1"}}
	idxPath := filepath.Join(dir, "indexes", "index_email_idx.json")
	f, _ := os.Create(idxPath)
	_ = json.NewEncoder(f).Encode(idxObj)
	f.Close()
	meta := fields.TableMeta{
		TableName: "users",
		Fields: []fields.FieldMeta{{Name: "id", Type: "string"}, {Name: "email", Type: "string"}},
		Indexes: []fields.IndexMeta{{Name: "email_idx", Type: "secondary", Fields: []string{"email"}, Unique: false}},
	}
	rep, err := CheckIndexConsistency(dir, meta, meta.Indexes[0])
	if err != nil {
		t.Fatalf("Fehler bei Konsistenzprüfung: %v", err)
	}
	if len(rep.MissingInIndex) != 1 || rep.MissingInIndex[0] != "id2" {
		t.Errorf("Fehlende Einträge im Index nicht erkannt: %+v", rep.MissingInIndex)
	}
	if len(rep.OrphanedInIndex) != 0 {
		t.Errorf("Orphaned falsch: %+v", rep.OrphanedInIndex)
	}
	if len(rep.Duplicates) != 0 {
		t.Errorf("Duplikate falsch: %+v", rep.Duplicates)
	}
}
