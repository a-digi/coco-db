package test

import (
	"os"
	"testing"
	"github.com/a-digi/coco-db/src/database"
)

func setupCreateTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_create_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownCreateTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestHandleCreateDatabase_Success(t *testing.T) {
	testDir := setupCreateTestDir(t)
	defer teardownCreateTestDir(testDir)
	dc := &database.DatabaseCreate{DataDir: testDir}
	resp := dc.HandleCreateDatabase("testdb")
	if resp == nil || !resp.Success || resp.HttpCode != 200 {
		t.Fatalf("Erwartet: Success true und HttpCode 200, erhalten: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
	if _, err := os.Stat(testDir + "/testdb"); err != nil {
		t.Errorf("Datenbankverzeichnis wurde nicht angelegt: %v", err)
	}
}

func TestHandleCreateDatabase_InvalidName(t *testing.T) {
	dc := &database.DatabaseCreate{DataDir: os.TempDir()}
	resp := dc.HandleCreateDatabase("!invalid")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_INVALID_NAME" || resp.HttpCode != 400 {
		t.Errorf("Ungültiger Name nicht korrekt erkannt: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
}

func TestHandleCreateDatabase_AlreadyExists(t *testing.T) {
	testDir := setupCreateTestDir(t)
	defer teardownCreateTestDir(testDir)
	os.MkdirAll(testDir+"/testdb", 0755)
	dc := &database.DatabaseCreate{DataDir: testDir}
	resp := dc.HandleCreateDatabase("testdb")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_EXISTS" || resp.HttpCode != 409 {
		t.Errorf("Doppelte Datenbank nicht korrekt erkannt: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
}
