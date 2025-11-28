package test

import (
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
	if resp == nil || !resp.Success {
		t.Fatalf("Erwartet: Success true, erhalten: %+v", resp)
	}
	dbs, ok := resp.Data.(map[string]interface{})["databases"].([]string)
	if !ok && resp.Data.(map[string]interface{})["databases"] != nil {
		// Fallback für []interface{}
		arr := resp.Data.(map[string]interface{})["databases"].([]interface{})
		if len(arr) != 0 {
			t.Errorf("Erwartet: 0 Datenbanken, erhalten: %d", len(arr))
		}
		return
	}
	if len(dbs) != 0 {
		t.Errorf("Erwartet: 0 Datenbanken, erhalten: %d", len(dbs))
	}
}

func TestHandleListDatabases_WithDatabases(t *testing.T) {
	testDir := setupListTestDir(t)
	defer teardownListTestDir(testDir)
	os.MkdirAll(testDir+"/db1", 0755)
	os.MkdirAll(testDir+"/db2", 0755)
	dl := &database.DatabaseList{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := dl.HandleListDatabases()
	if resp == nil || !resp.Success {
		t.Fatalf("Erwartet: Success true, erhalten: %+v", resp)
	}
	dbsIface := resp.Data.(map[string]interface{})["databases"]
	dbs := []string{}
	switch v := dbsIface.(type) {
	case []interface{}:
		for _, entry := range v {
			if s, ok := entry.(string); ok {
				dbs = append(dbs, s)
			}
		}
	case []string:
		dbs = v
	}
	if len(dbs) != 2 {
		t.Errorf("Erwartet: 2 Datenbanken, erhalten: %d", len(dbs))
	}
}

