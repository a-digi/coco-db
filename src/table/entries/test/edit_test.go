package entries_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/index"
	"github.com/a-digi/coco-db/src/table/entries"
	"github.com/a-digi/coco-db/src/table/entries/testmeta"
)

func TestEditEntry_UpdatesMemoryIndex(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "editentry_test_memidx")
	if err != nil {
		t.Fatalf("TempDir error: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbName := "testdb"
	tableName := "testtable"
	entryId := "testid"
	origEntry := map[string]interface{}{"foo": "bar"}

	// meta.json mit Index "foo" anlegen
	tableDir := filepath.Join(tmpDir, dbName, tableName)
	_ = os.MkdirAll(tableDir, 0755)
	meta := testmeta.TestMeta{
		Indexes: []struct {
			Name   string   `json:"name"`
			Fields []string `json:"fields"`
		}{
			{Name: "foo", Fields: []string{"foo"}},
		},
		TotalEntries: 0,
	}
	metaPath := filepath.Join(tableDir, "meta.json")
	metaFile, _ := os.Create(metaPath)
	_ = json.NewEncoder(metaFile).Encode(meta)
	metaFile.Close()

	reg := index.GetRegistry()
	idxName := "foo"
	reg.Set(dbName+"."+tableName+"."+idxName, index.IndexData{})

	// Insert vor Edit
	_ = entries.InsertEntry(dbName, tableName, origEntry, tmpDir, nil)

	entriesEditor := &entries.EntryEditor{DataDir: tmpDir, Logger: nil}
	entriesEditor.EditEntry(dbName, tableName, entryId, origEntry)
	idxKey := "bar"
	idx, ok := reg.Get(dbName + "." + tableName + "." + idxName)
	if !ok {
		t.Fatalf("Index nicht im Speicher gefunden")
	}
	_, ok = idx[idxKey].([]string)
	if !ok {
		t.Fatalf("Key nicht im Speicherindex gefunden")
	}
}
