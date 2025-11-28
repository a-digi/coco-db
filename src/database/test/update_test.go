package test

import (
	"os"
	"testing"
	"github.com/a-digi/coco-db/src/database"
)

func setupUpdateTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "cocodb_update_test_")
	if err != nil {
		t.Fatalf("TempDir Fehler: %v", err)
	}
	return dir
}

func teardownUpdateTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestHandleUpdateDatabase_Success(t *testing.T) {
	testDir := setupUpdateTestDir(t)
	defer teardownUpdateTestDir(testDir)
	oldName := "olddb"
	newName := "newdb"
	os.MkdirAll(testDir+"/"+oldName, 0755)
	du := &database.DatabaseUpdate{SataDir: testDir}
	resp := du.HandleUpdateDatabase(oldName, newName)
	if resp == nil || !resp.Success || resp.HttpCode != 200 {
		t.Fatalf("Erwartet: Success true und HttpCode 200, erhalten: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
	if _, err := os.Stat(testDir + "/" + newName); err != nil {
		t.Errorf("Neuer Datenbankordner wurde nicht angelegt: %v", err)
	}
	if _, err := os.Stat(testDir + "/" + oldName); err == nil {
		t.Errorf("Alter Datenbankordner existiert noch")
	}
}

func TestHandleUpdateDatabase_InvalidOldName(t *testing.T) {
	du := &database.DatabaseUpdate{SataDir: os.TempDir()}
	resp := du.HandleUpdateDatabase("!invalid", "newdb")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_INVALID_NAME" || resp.HttpCode != 400 {
		t.Errorf("Ungültiger alter Name nicht korrekt erkannt: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
}

func TestHandleUpdateDatabase_InvalidNewName(t *testing.T) {
	du := &database.DatabaseUpdate{SataDir: os.TempDir()}
	resp := du.HandleUpdateDatabase("olddb", "!invalid")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_INVALID_NAME" || resp.HttpCode != 400 {
		t.Errorf("Ungültiger neuer Name nicht korrekt erkannt: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
}

func TestHandleUpdateDatabase_SameName(t *testing.T) {
	du := &database.DatabaseUpdate{SataDir: os.TempDir()}
	resp := du.HandleUpdateDatabase("olddb", "olddb")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_SAME_NAME" || resp.HttpCode != 400 {
		t.Errorf("Gleicher Name nicht korrekt erkannt: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
}

func TestHandleUpdateDatabase_NotFound(t *testing.T) {
	testDir := setupUpdateTestDir(t)
	defer teardownUpdateTestDir(testDir)
	du := &database.DatabaseUpdate{SataDir: testDir}
	resp := du.HandleUpdateDatabase("notfound", "newdb")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_NOT_FOUND" || resp.HttpCode != 404 {
		t.Errorf("Nicht vorhandene Datenbank nicht korrekt erkannt: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
}

func TestHandleUpdateDatabase_AlreadyExists(t *testing.T) {
	testDir := setupUpdateTestDir(t)
	defer teardownUpdateTestDir(testDir)
	os.MkdirAll(testDir+"/olddb", 0755)
	os.MkdirAll(testDir+"/newdb", 0755)
	du := &database.DatabaseUpdate{SataDir: testDir}
	resp := du.HandleUpdateDatabase("olddb", "newdb")
	if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != "ERR_DB_EXISTS" || resp.HttpCode != 409 {
		t.Errorf("Doppelte Datenbank nicht korrekt erkannt: %+v", resp)
	}
	if resp.ExecutionTime == "" {
		t.Errorf("ExecutionTime fehlt")
	}
}
