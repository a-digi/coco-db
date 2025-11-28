package test

import (
	"os"
	"testing"
	"github.com/a-digi/coco-db/src/database"
	"github.com/a-digi/coco-db/src/logger"
)

func setupDeleteTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_delete_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownDeleteTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestHandleDeleteDatabase_Success(t *testing.T) {
	testDir := setupDeleteTestDir(t)
	defer teardownDeleteTestDir(testDir)
	dbName := "testdb"
	dbPath := testDir + "/" + dbName
	os.MkdirAll(dbPath, 0755)
	dd := &database.DatabaseDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := dd.HandleDeleteDatabase(dbName)
	if resp == nil || !resp.Success {
		t.Fatalf("Erwartet: Success true, erhalten: %+v", resp)
	}
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Errorf("Datenbankverzeichnis wurde nicht gelöscht: %v", err)
	}
}

func TestHandleDeleteDatabase_NotFound(t *testing.T) {
	testDir := setupDeleteTestDir(t)
	defer teardownDeleteTestDir(testDir)
	dbName := "notfound"
	dd := &database.DatabaseDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := dd.HandleDeleteDatabase(dbName)
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_NOT_FOUND" {
		t.Errorf("Nicht vorhandene Datenbank nicht korrekt erkannt: %+v", resp)
	}
}

func TestHandleDeleteDatabase_InvalidName(t *testing.T) {
	testDir := setupDeleteTestDir(t)
	defer teardownDeleteTestDir(testDir)
	dd := &database.DatabaseDelete{DataDir: testDir, Logger: &logger.NoopLogger{}}
	resp := dd.HandleDeleteDatabase("!invalid")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_INVALID_NAME" {
		t.Errorf("Ungültiger Name nicht korrekt erkannt: %+v", resp)
	}
}
