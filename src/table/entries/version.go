package entries

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// VersionEntry beschreibt eine Version eines Eintrags
//
type VersionEntry struct {
	EntryID      string `json:"entryId"`
	CreatedAt    string `json:"created_at"`
	VersionedAt  string `json:"versioned_at"`
	VersionNumber int    `json:"version_number"`
}

// Versioning kapselt die Versionierungslogik für Einträge
//
type Versioning struct {
	Dir string // Pfad zu <DataDir>/<dbName>/<tableName>/entries/<entryId>
}

func (v *Versioning) versionFilePath() string {
	return filepath.Join(v.Dir, "version.json")
}

func (v *Versioning) LoadVersions() ([]VersionEntry, error) {
	path := v.versionFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []VersionEntry{}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var versions []VersionEntry
	if err := json.NewDecoder(f).Decode(&versions); err != nil {
		return nil, err
	}
	return versions, nil
}

func (v *Versioning) SaveVersions(versions []VersionEntry) error {
	path := v.versionFilePath()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(versions)
}

// Versioniere den aktuellen Eintrag, falls vorhanden, und gib die nächste VersionNumber zurück
func (v *Versioning) VersionActiveEntry() (int, error) {
	versions, err := v.LoadVersions()
	if err != nil {
		return 0, err
	}
	var maxVersion int
	for _, ve := range versions {
		if ve.VersionNumber > maxVersion {
			maxVersion = ve.VersionNumber
		}
	}
	// Finde aktiven Eintrag (versioned_at leer)
	for i, ve := range versions {
		if ve.VersionedAt == "" {
			// Versioniere: setze versioned_at
			versions[i].VersionedAt = time.Now().UTC().Format(time.RFC3339)
			// Datei verschieben
			versionDir := filepath.Join(v.Dir, "version")
			if err := os.MkdirAll(versionDir, 0755); err != nil {
				return 0, err
			}
			oldPath := filepath.Join(v.Dir, ve.EntryID+".json")
			newPath := filepath.Join(versionDir, fmt.Sprintf("%d.json", ve.VersionNumber))
			if err := os.Rename(oldPath, newPath); err != nil {
				return 0, err
			}
			break
		}
	}
	// Sortiere nach VersionNumber
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].VersionNumber < versions[j].VersionNumber
	})
	if err := v.SaveVersions(versions); err != nil {
		return 0, err
	}
	return maxVersion + 1, nil
}

// Füge eine neue Version als aktuell hinzu
func (v *Versioning) AddNewVersion(entryId string, versionNumber int) error {
	versions, err := v.LoadVersions()
	if err != nil {
		return err
	}

	createdAt := time.Now().UTC().Format(time.RFC3339)
	ve := VersionEntry{
		EntryID:      entryId,
		CreatedAt:    createdAt,
		VersionedAt:  "",
		VersionNumber: versionNumber,
	}
	versions = append(versions, ve)
	// Sortiere nach VersionNumber
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].VersionNumber < versions[j].VersionNumber
	})
	return v.SaveVersions(versions)
}

// Entferne MarkDeleted (Soft Delete) und implementiere stattdessen Hard Delete
// Hard Delete: Lösche alle Dateien und den Ordner des Eintrags inklusive aller Versionen und Metadaten
func (v *Versioning) HardDelete() error {
	// Lösche das gesamte Verzeichnis (inkl. version.json, Versionen, JSON-Dateien)
	return os.RemoveAll(v.Dir)
}

// VersionEntryWithDelete erweitert VersionEntry um deleted_at
// (für Soft Delete, falls benötigt)
type VersionEntryWithDelete struct {
	VersionEntry
	DeletedAt string `json:"deleted_at,omitempty"`
}
