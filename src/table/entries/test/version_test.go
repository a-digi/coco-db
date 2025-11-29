package entries_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	entries "github.com/a-digi/coco-db/src/table/entries"
)

func TestVersioning_AddAndVersionEntry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "versioning_test")
	if err != nil {
		t.Fatalf("TempDir error: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	entryId := "testentry"
	entryDir := filepath.Join(tmpDir, entryId)
	if err := os.MkdirAll(entryDir, 0755); err != nil {
		t.Fatalf("MkdirAll error: %v", err)
	}
	v := entries.Versioning{Dir: entryDir}

	// Add first version
	err = v.AddNewVersion(entryId, 1)
	if err != nil {
		t.Fatalf("AddNewVersion error: %v", err)
	}

	// Prüfe version.json
	versionPath := filepath.Join(entryDir, "version.json")
	data, err := os.ReadFile(versionPath)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	var versions []entries.VersionEntry
	if err := json.Unmarshal(data, &versions); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if len(versions) != 1 {
		t.Errorf("expected 1 version, got %d", len(versions))
	}
	if versions[0].VersionNumber != 1 {
		t.Errorf("expected version_number 1, got %d", versions[0].VersionNumber)
	}
	if versions[0].VersionedAt != "" {
		t.Errorf("expected versioned_at to be empty, got %s", versions[0].VersionedAt)
	}

	// Simuliere Datei für Version 1
	f, err := os.Create(filepath.Join(entryDir, entryId+".json"))
	if err != nil {
		t.Fatalf("Create file error: %v", err)
	}
	f.Close()

	// Warte eine Sekunde, damit versioned_at unterschiedlich ist
	time.Sleep(1 * time.Second)

	// Versioniere
	nextVersion, err := v.VersionActiveEntry()
	if err != nil {
		t.Fatalf("VersionActiveEntry error: %v", err)
	}
	if nextVersion != 2 {
		t.Errorf("expected nextVersion 2, got %d", nextVersion)
	}
	// Prüfe, ob Datei verschoben wurde
	versionedFile := filepath.Join(entryDir, "version", "1.json")
	if _, err := os.Stat(versionedFile); os.IsNotExist(err) {
		t.Errorf("versioned file not found: %s", versionedFile)
	}
	// Prüfe versioned_at gesetzt
	data, err = os.ReadFile(versionPath)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if err := json.Unmarshal(data, &versions); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if versions[0].VersionedAt == "" {
		t.Errorf("expected versioned_at to be set, got empty")
	}
}

