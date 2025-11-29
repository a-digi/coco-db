package index

import (
	"encoding/json"
	"os"
	"path/filepath"
	"fmt"
	fields "github.com/a-digi/coco-db/src/table/fields"
)

// ConsistencyReport enthält die Ergebnisse der Konsistenzprüfung
// - MissingInIndex: Einträge, die in den Daten existieren, aber nicht im Index
// - OrphanedInIndex: Einträge, die im Index existieren, aber nicht in den Daten
// - Duplicates: Einträge, die mehrfach im Index stehen (bei unique)
type ConsistencyReport struct {
	IndexName      string
	MissingInIndex []string // entryIds
	OrphanedInIndex []string // entryIds
	Duplicates     []string // entryIds
}

// CheckIndexConsistency prüft für einen Single-Field-Index die Konsistenz zwischen Daten und Indexdatei
func CheckIndexConsistency(tableDir string, meta fields.TableMeta, idxMeta fields.IndexMeta) (*ConsistencyReport, error) {
	entriesDir := filepath.Join(tableDir, "entries")
	idxPath := filepath.Join(tableDir, "index_"+idxMeta.Name+".json")
	// 1. Alle EntryIDs aus Daten sammeln
	entryIds := map[string]string{} // entryId -> key
	files, err := os.ReadDir(entriesDir)
	if err != nil {
		return nil, err
	}
	for _, entryFolder := range files {
		if !entryFolder.IsDir() { continue }
		entryId := entryFolder.Name()
		entryFile := filepath.Join(entriesDir, entryId, entryId+".json")
		data, err := os.ReadFile(entryFile)
		if err != nil { continue }
		var entry map[string]interface{}
		if err := json.Unmarshal(data, &entry); err != nil { continue }
		if len(idxMeta.Fields) == 1 {
			key, ok := entry[idxMeta.Fields[0]]
			if ok {
				k, ok := key.(string)
				if ok {
					entryIds[entryId] = k
				}
			}
		}
	}
	// 2. Alle EntryIDs aus Indexdatei sammeln
	idxObj := map[string][]string{}
	if idxData, err := os.ReadFile(idxPath); err == nil {
		_ = json.Unmarshal(idxData, &idxObj)
	}
	indexEntryIds := map[string]string{} // entryId -> key
	duplicates := map[string]bool{}
	for key, ids := range idxObj {
		seen := map[string]bool{}
		for _, id := range ids {
			if seen[id] {
				duplicates[id] = true
			} else {
				seen[id] = true
			}
			indexEntryIds[id] = key
		}
	}
	// 3. Vergleiche
	missing := []string{}
	for id := range entryIds {
		if _, ok := indexEntryIds[id]; !ok {
			missing = append(missing, id)
		}
	}
	orphaned := []string{}
	for id := range indexEntryIds {
		if _, ok := entryIds[id]; !ok {
			orphaned = append(orphaned, id)
		}
	}
	dupes := []string{}
	for id := range duplicates {
		dupes = append(dupes, id)
	}
	return &ConsistencyReport{
		IndexName: idxMeta.Name,
		MissingInIndex: missing,
		OrphanedInIndex: orphaned,
		Duplicates: dupes,
	}, nil
}

// CheckAllIndexesConsistency prüft alle Indizes einer Tabelle und gibt einen Report je Index zurück
func CheckAllIndexesConsistency(tableDir string, meta fields.TableMeta) ([]*ConsistencyReport, error) {
	reports := []*ConsistencyReport{}
	for _, idxMeta := range meta.Indexes {
		if len(idxMeta.Fields) == 1 {
			rep, err := CheckIndexConsistency(tableDir, meta, idxMeta)
			if err != nil {
				return nil, fmt.Errorf("Fehler bei Index %s: %v", idxMeta.Name, err)
			}
			reports = append(reports, rep)
		}
	}
	return reports, nil
}

