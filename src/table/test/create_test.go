package table_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/table"
	"github.com/a-digi/coco-db/src/table/fields"
)

func setupTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_create_index_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownTestDir(dir string) {
	_ = os.RemoveAll(dir)
}

func TestHandleCreateTable_WithIndex_Success(t *testing.T) {
	t.Log("TestHandleCreateTable_WithIndex_Success läuft")
	testDir := setupTestDir(t)
	defer teardownTestDir(testDir)
	creator := &table.TableCreator{
		DataDir: testDir,
		Logger:  &logger.NoopLogger{},
	}
	meta := fields.TableMeta{
		TableName: "users",
		Fields:    []fields.FieldMeta{{Name: "id", Type: "string", Required: true}},
		Indexes: []fields.IndexMeta{
			{Name: "primary_id", Type: "primary", Fields: []string{"id"}, Unique: true},
		},
	}
	resp := creator.HandleCreateTable("testdb", meta)
	if resp == nil || !resp.Success {
		t.Fatalf("Erwartet: Success true, erhalten: %+v", resp)
	}
	// Prüfe, ob meta.json existiert und Index korrekt gespeichert ist
	metaPath := filepath.Join(testDir, "testdb", "users", "meta.json")
	content, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("meta.json wurde nicht angelegt: %v", err)
	}
	var loadedMeta fields.TableMeta
	if err := json.Unmarshal(content, &loadedMeta); err != nil {
		t.Fatalf("meta.json nicht parsebar: %v", err)
	}
	if len(loadedMeta.Indexes) != 1 {
		t.Errorf("Index nicht korrekt gespeichert: %+v", loadedMeta.Indexes)
	}
	if loadedMeta.Indexes[0].Name != "primary_id" || !loadedMeta.Indexes[0].Unique {
		t.Errorf("Indexdaten stimmen nicht: %+v", loadedMeta.Indexes[0])
	}
}
