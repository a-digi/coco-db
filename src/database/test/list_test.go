package test

import (
	"encoding/json"
	"os"
	"testing"
	"github.com/a-digi/coco-db/src/database"
	"github.com/a-digi/coco-db/src/logger"
)

func setupListTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_list_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownListTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestHandleListDatabases_Empty(t *testing.T) {
	testDir := setupListTestDir(t)
	defer teardownListTestDir(testDir)
	dl := &database.DatabaseList{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := dl.HandleListDatabases()
	if resp == nil || !resp.Success || resp.HttpCode != 200 {
		t.Fatalf("Erwartet: Success true und HttpCode 200, erhalten: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
	dbsIface := resp.Data.(map[string]interface{})["databases"]
	dbs := []database.DatabaseMeta{}
	switch v := dbsIface.(type) {
	case []interface{}:
		for _, entry := range v {
			if m, ok := entry.(map[string]interface{}); ok {
				if name, ok := m["name"].(string); ok {
					dbs = append(dbs, database.DatabaseMeta{Name: name})
				}
			}
		}
	case []database.DatabaseMeta:
		dbs = v
	}
	if len(dbs) != 0 {
		t.Errorf("Erwartet: 0 Datenbanken, erhalten: %d", len(dbs))
	}
}

func TestHandleListDatabases_WithDatabases(t *testing.T) {
	testDir := setupListTestDir(t)
	defer teardownListTestDir(testDir)
	// Simuliere tables.json mit zwei Datenbanken
	dbs := []database.DatabaseMeta{{Name: "db1"}, {Name: "db2"}}
	tablesPath := testDir + "/tables.json"
	f, err := os.Create(tablesPath)
	if err != nil {
		t.Fatalf("Fehler beim Anlegen von tables.json: %v", err)
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(dbs)
	dl := &database.DatabaseList{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := dl.HandleListDatabases()
	if resp == nil || !resp.Success || resp.HttpCode != 200 {
		t.Fatalf("Erwartet: Success true und HttpCode 200, erhalten: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
	dbsIface := resp.Data.(map[string]interface{})["databases"]
	names := []string{}
	switch v := dbsIface.(type) {
	case []interface{}:
		for _, entry := range v {
			if m, ok := entry.(map[string]interface{}); ok {
				if name, ok := m["name"].(string); ok {
					names = append(names, name)
				}
			}
		}
	case []database.DatabaseMeta:
		for _, m := range v {
			names = append(names, m.Name)
		}
	}
	if len(names) != 2 {
		t.Errorf("Erwartet: 2 Datenbanken, erhalten: %d", len(names))
	}
}
