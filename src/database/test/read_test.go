package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/a-digi/coco-db/src/database"
)

func setupTestDataDirRead(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_read_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownTestDataDirRead(dir string) {
	os.RemoveAll(dir)
}

func TestListDatabases_Empty(t *testing.T) {
	dataDir := setupTestDataDirRead(t)
	defer teardownTestDataDirRead(dataDir)

	r := httptest.NewRequest(http.MethodGet, "/api/databases", nil)
	w := httptest.NewRecorder()
	resp := database.HandleListDatabases(w, r, dataDir)
	if !resp.Success {
		t.Errorf("Datenbanken auflisten fehlgeschlagen: %+v", resp)
	}
}

