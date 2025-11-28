package database

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestDataDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "cocodb_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownTestDataDir(dir string) {
	os.RemoveAll(dir)
}

func TestCreateListDeleteDatabase(t *testing.T) {
	dataDir := setupTestDataDir(t)
	defer teardownTestDataDir(dataDir)

	// 1. Datenbank anlegen
	body := bytes.NewBufferString(`{"name":"testdb"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/databases", body)
	w := httptest.NewRecorder()
	resp := handleCreateDatabase(w, r, dataDir)
	if !resp.Success {
		t.Errorf("Datenbank anlegen fehlgeschlagen: %+v", resp)
	}

	// 2. Datenbank erneut anlegen (sollte Fehler geben)
	body2 := bytes.NewBufferString(`{"name":"testdb"}`)
	r2 := httptest.NewRequest(http.MethodPost, "/api/databases", body2)
	w2 := httptest.NewRecorder()
	resp2 := handleCreateDatabase(w2, r2, dataDir)
	if resp2.Success || resp2.Error == nil || resp2.Error.Code != "ERR_DB_EXISTS" {
		t.Errorf("Doppelte Datenbank nicht korrekt erkannt: %+v", resp2)
	}

	// 3. Datenbanken auflisten
	r3 := httptest.NewRequest(http.MethodGet, "/api/databases", nil)
	w3 := httptest.NewRecorder()
	resp3 := handleListDatabases(w3, r3, dataDir)
	if !resp3.Success {
		t.Errorf("Datenbanken auflisten fehlgeschlagen: %+v", resp3)
	}
	var out struct{ Databases []string `json:"databases"` }
	b, _ := json.Marshal(resp3.Data)
	json.Unmarshal(b, &out)
	if len(out.Databases) != 1 || out.Databases[0] != "testdb" {
		t.Errorf("Datenbank nicht korrekt gelistet: %+v", out.Databases)
	}

	// 4. Datenbank löschen
	r4 := httptest.NewRequest(http.MethodGet, "/api/databases/delete/testdb", nil)
	w4 := httptest.NewRecorder()
	resp4 := handleDeleteDatabase(w4, r4, dataDir)
	if !resp4.Success {
		t.Errorf("Datenbank löschen fehlgeschlagen: %+v", resp4)
	}

	// 5. Datenbank löschen (nicht gefunden)
	r5 := httptest.NewRequest(http.MethodGet, "/api/databases/delete/testdb", nil)
	w5 := httptest.NewRecorder()
	resp5 := handleDeleteDatabase(w5, r5, dataDir)
	if resp5.Success || resp5.Error == nil || resp5.Error.Code != "ERR_DB_NOT_FOUND" {
		t.Errorf("Nicht vorhandene Datenbank nicht korrekt erkannt: %+v", resp5)
	}
}

func TestInvalidDatabaseName(t *testing.T) {
	dataDir := setupTestDataDir(t)
	defer teardownTestDataDir(dataDir)

	body := bytes.NewBufferString(`{"name":"../hack"}`)
	r := httptest.NewRequest(http.MethodPost, "/api/databases", body)
	w := httptest.NewRecorder()
	resp := handleCreateDatabase(w, r, dataDir)
	if resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_INVALID_NAME" {
		t.Errorf("Ungültiger Name nicht korrekt erkannt: %+v", resp)
	}
}
