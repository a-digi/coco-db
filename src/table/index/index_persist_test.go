package index

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBTreeIndex_SaveAndLoadToFile(t *testing.T) {
	dir := t.TempDir()
	meta := IndexMeta{Name: "email", Fields: []string{"email"}, Type: IndexTypeBTree}
	idx := NewBTreeIndex(meta)
	// Simuliere einen Index als Map
	stub := map[string][]string{"foo@bar.de": {"id1"}, "bar@baz.de": {"id2", "id3"}}
	idx.fromMap(stub)
	if err := idx.SaveToFile(dir); err != nil {
		t.Fatalf("Fehler beim Speichern: %v", err)
	}
	// Lade in neuen Index
	idx2 := NewBTreeIndex(meta)
	if err := idx2.LoadFromFile(dir); err != nil {
		t.Fatalf("Fehler beim Laden: %v", err)
	}
	m := idx2.toMap()
	if len(m) != 2 || len(m["foo@bar.de"]) != 1 || len(m["bar@baz.de"]) != 2 {
		t.Errorf("Indexdaten nach Laden stimmen nicht: %+v", m)
	}
	// Prüfe, ob Datei existiert
	indexPath := filepath.Join(dir, "index_email.json")
	if _, err := os.Stat(indexPath); err != nil {
		t.Errorf("Indexdatei nicht gefunden: %v", err)
	}
}

