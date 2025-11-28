package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/a-digi/coco-db/src/database"
)

func setupTestDataDirDelete(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_delete_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownTestDataDirDelete(dir string) {
	os.RemoveAll(dir)
}

func TestDeleteDatabase_Success(t *testing.T) {
	dataDir := setupTestDataDirDelete(t)
	defer teardownTestDataDirDelete(dataDir)

	// Erst anlegen
	body := bytes.NewBufferString(`{"name":"testdb"}`)
	rCreate := httptest.NewRequest(http.MethodGet, "/api/databases/create/", body)
	wCreate := httptest.NewRecorder()
	database.HandleCreateDatabase(wCreate, rCreate, dataDir)

	// Dann löschen
	r := httptest.NewRequest(http.MethodGet, "/api/databases/delete/testdb", nil)
	w := httptest.NewRecorder()
	resp := database.HandleDeleteDatabase(w, r, dataDir)
	if !resp.Success {
		t.Errorf("Datenbank löschen fehlgeschlagen: %+v", resp)
		if resp.Error != nil {
			t.Logf("Fehlercode: %s, Nachricht: %s", resp.Error.Code, resp.Error.Message)
		}
	}
}

func TestDeleteDatabase_NotFound(t *testing.T) {
	dataDir := setupTestDataDirDelete(t)
	defer teardownTestDataDirDelete(dataDir)

	r := httptest.NewRequest(http.MethodGet, "/api/databases/delete/notfounddb", nil)
	w := httptest.NewRecorder()
	resp := database.HandleDeleteDatabase(w, r, dataDir)
	if resp.Success {
		t.Errorf("Nicht vorhandene Datenbank als Erfolg behandelt: %+v", resp)
	}
	if resp.Error == nil {
		t.Errorf("Fehlerobjekt fehlt: %+v", resp)
	} else if resp.Error.Code != "ERR_DB_NOT_FOUND" {
		t.Errorf("Falscher Fehlercode: %s, Nachricht: %s", resp.Error.Code, resp.Error.Message)
	}
}
