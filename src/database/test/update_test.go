package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/a-digi/coco-db/src/database"
)

func setupTestDataDirUpdate(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_update_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownTestDataDirUpdate(dir string) {
	os.RemoveAll(dir)
}

func TestUpdateDatabase_SameName(t *testing.T) {
	dataDir := setupTestDataDirUpdate(t)
	defer teardownTestDataDirUpdate(dataDir)

	// Erst anlegen
	body := bytes.NewBufferString(`{"name":"testdb"}`)
	rCreate := httptest.NewRequest(http.MethodPost, "/api/databases", body)
	wCreate := httptest.NewRecorder()
	database.HandleCreateDatabase(wCreate, rCreate, dataDir)

	// Dann Update mit gleichem Namen
	updateBody := bytes.NewBufferString(`{"newName":"testdb"}`)
	r := httptest.NewRequest(http.MethodPut, "/api/databases/testdb", updateBody)
	w := httptest.NewRecorder()
	resp := database.HandleUpdateDatabase(w, r, dataDir)
	if resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_SAME_NAME" {
		t.Errorf("Update mit identischem Namen nicht korrekt erkannt: %+v", resp)
	}
}

func TestUpdateDatabase_AlreadyExists(t *testing.T) {
	dataDir := setupTestDataDirUpdate(t)
	defer teardownTestDataDirUpdate(dataDir)

	// Zwei Datenbanken anlegen
	body1 := bytes.NewBufferString(`{"name":"db1"}`)
	rCreate1 := httptest.NewRequest(http.MethodPost, "/api/databases", body1)
	wCreate1 := httptest.NewRecorder()
	database.HandleCreateDatabase(wCreate1, rCreate1, dataDir)

	body2 := bytes.NewBufferString(`{"name":"db2"}`)
	rCreate2 := httptest.NewRequest(http.MethodPost, "/api/databases", body2)
	wCreate2 := httptest.NewRecorder()
	database.HandleCreateDatabase(wCreate2, rCreate2, dataDir)

	// Update db1 -> db2 (sollte Fehler geben)
	updateBody := bytes.NewBufferString(`{"newName":"db2"}`)
	r := httptest.NewRequest(http.MethodPut, "/api/databases/db1", updateBody)
	w := httptest.NewRecorder()
	resp := database.HandleUpdateDatabase(w, r, dataDir)
	if resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_EXISTS" {
		t.Errorf("Update auf bestehenden Namen nicht korrekt erkannt: %+v", resp)
	}
}

func TestUpdateDatabase_Success(t *testing.T) {
	dataDir := setupTestDataDirUpdate(t)
	defer teardownTestDataDirUpdate(dataDir)

	// Datenbank anlegen
	body := bytes.NewBufferString(`{"name":"db1"}`)
	rCreate := httptest.NewRequest(http.MethodPost, "/api/databases", body)
	wCreate := httptest.NewRecorder()
	database.HandleCreateDatabase(wCreate, rCreate, dataDir)

	// Update db1 -> db2
	updateBody := bytes.NewBufferString(`{"newName":"db2"}`)
	r := httptest.NewRequest(http.MethodPut, "/api/databases/db1", updateBody)
	w := httptest.NewRecorder()
	resp := database.HandleUpdateDatabase(w, r, dataDir)
	if !resp.Success {
		t.Errorf("Update fehlgeschlagen: %+v", resp)
	}
}
