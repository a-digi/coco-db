package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/a-digi/coco-db/src/database"
)

func setupTestDataDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_create_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownTestDataDir(dir string) {
	os.RemoveAll(dir)
}

func TestCreateDatabase_Success(t *testing.T) {
	dataDir := setupTestDataDir(t)
	defer teardownTestDataDir(dataDir)

	body := bytes.NewBufferString(`{"name":"testdb"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/databases", body)
	w := httptest.NewRecorder()
	resp := database.HandleCreateDatabase(w, r, dataDir)
	if !resp.Success {
		t.Errorf("Datenbank anlegen fehlgeschlagen: %+v", resp)
	}
}

func TestCreateDatabase_InvalidName(t *testing.T) {
	dataDir := setupTestDataDir(t)
	defer teardownTestDataDir(dataDir)

	body := bytes.NewBufferString(`{"name":"../hack"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/databases", body)
	w := httptest.NewRecorder()
	resp := database.HandleCreateDatabase(w, r, dataDir)
	if resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_INVALID_NAME" {
		t.Errorf("Ungültiger Name nicht korrekt erkannt: %+v", resp)
	}
}

func TestCreateDatabase_AlreadyExists(t *testing.T) {
	dataDir := setupTestDataDir(t)
	defer teardownTestDataDir(dataDir)

	body := bytes.NewBufferString(`{"name":"testdb"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/databases", body)
	w := httptest.NewRecorder()
	database.HandleCreateDatabase(w, r, dataDir)

	body2 := bytes.NewBufferString(`{"name":"testdb"}`)
	r2 := httptest.NewRequest(http.MethodPost, "/api/databases", body2)
	w2 := httptest.NewRecorder()
	resp2 := database.HandleCreateDatabase(w2, r2, dataDir)
	if resp2.Success || resp2.Error == nil || resp2.Error.Code != "ERR_DB_EXISTS" {
		t.Errorf("Doppelte Datenbank nicht korrekt erkannt: %+v", resp2)
	}
}
