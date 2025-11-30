package index

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestIndexRegistry_SetGet(t *testing.T) {
	reg := GetRegistry()
	key := "testdb.testtable.testindex"
	data := IndexData{"foo": "bar"}
	reg.Set(key, data)

	got, ok := reg.Get(key)
	if !ok {
		t.Fatalf("Index %s not found", key)
	}
	if got["foo"] != "bar" {
		t.Errorf("Expected foo=bar, got %v", got["foo"])
	}
}

func TestLoadAllIndexes(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "index_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbDir := filepath.Join(tmpDir, "testdb", "testtable")
	os.MkdirAll(dbDir, 0755)
	idxPath := filepath.Join(dbDir, "index_test.json")
	idxContent := []byte(`{"foo": "baz"}`)
	if err := ioutil.WriteFile(idxPath, idxContent, 0644); err != nil {
		t.Fatal(err)
	}

	registry = nil // reset singleton
	if err := LoadAllIndexes(tmpDir); err != nil {
		t.Fatalf("LoadAllIndexes failed: %v", err)
	}

	reg := GetRegistry()
	key := "testdb.testtable.test"
	got, ok := reg.Get(key)
	if !ok {
		t.Fatalf("Index %s not found after loading", key)
	}
	if got["foo"] != "baz" {
		t.Errorf("Expected foo=baz, got %v", got["foo"])
	}
}

