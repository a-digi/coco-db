package query

import (
	"encoding/json"
	"os"
	"path/filepath"
	"github.com/a-digi/coco-db/src/table/fields"
)

// Speicheroptimierte Filter-Engine: Nur Indexdaten im Speicher, sonst sequentieller Dateiscan
// Gibt die gefilterten Einträge als Array von map[string]interface{} zurück
func FilterEngine(dataDir, dbName, tableName string, query *Query, meta *fields.TableMeta) ([]map[string]interface{}, error) {
	// 1. Indexfelder und Nicht-Indexfelder trennen
	indexedFields := map[string]fields.IndexMeta{}
	nonIndexedFields := map[string]struct{}{}
	for _, idx := range meta.Indexes {
		for _, f := range idx.Fields {
			indexedFields[f] = idx
		}
	}
	for f := range query.Filter {
		if _, ok := indexedFields[f]; !ok {
			nonIndexedFields[f] = struct{}{}
		}
	}

	// 2. Indexnutzung: IDs für alle Indexfilter sammeln (Schnittmenge)
	var idSets [][]string
	for f, idxMeta := range indexedFields {
		if cond, ok := query.Filter[f]; ok {
			// Index laden (Stub: Annahme BTree)
			idxPath := filepath.Join(dataDir, dbName, tableName, "indexes", idxMeta.Name+".json")
			ids, err := loadIDsFromIndex(idxPath, cond)
			if err != nil {
				return nil, err
			}
			idSets = append(idSets, ids)
		}
	}
	ids := intersectIDSets(idSets)

	// 3. Sequentieller Scan für nicht indizierte Filter
	entriesDir := filepath.Join(dataDir, dbName, tableName, "entries")
	var result []map[string]interface{}
	if len(ids) > 0 {
		// Nur Einträge mit diesen IDs laden
		for _, id := range ids {
			entry, err := loadEntry(entriesDir, id)
			if err == nil && matchesAllFiltersEngine(entry, query.Filter) {
				result = append(result, entry)
			}
		}
	} else {
		// Kein Indexfilter: Alle Einträge sequenziell prüfen (ohne alles in den Speicher zu laden)
		files, _ := os.ReadDir(entriesDir)
		for _, f := range files {
			if f.IsDir() {
				id := f.Name()
				entry, err := loadEntry(entriesDir, id)
				if err == nil && matchesAllFiltersEngine(entry, query.Filter) {
					result = append(result, entry)
				}
			}
		}
	}
	return result, nil
}

// matchesAllFiltersEngine prüft, ob ein Eintrag alle Filterbedingungen erfüllt (AND-Logik)
func matchesAllFiltersEngine(entry map[string]interface{}, filter map[string]interface{}) bool {
	for field, cond := range filter {
		val, ok := entry[field]
		if !ok {
			return false
		}
		switch c := cond.(type) {
		case map[string]interface{}:
			for op, opVal := range c {
				fn, found := operatorFuncs[op]
				if !found {
					continue
				}
				if !fn(val, opVal) {
					return false
				}
			}
			continue
		default:
			if !isEqual(val, c) {
				return false
			}
		}
	}
	return true
}

// Hilfsfunktion: IDs aus Indexdatei laden (Stub: alle IDs zurückgeben, echte Filterung TODO)
func loadIDsFromIndex(idxPath string, cond interface{}) ([]string, error) {
	b, err := os.ReadFile(idxPath)
	if err != nil {
		return nil, err
	}
	var ids []string
	_ = json.Unmarshal(b, &ids) // Annahme: Flat-Array, TODO: echte Indexlogik
	return ids, nil
}

// Hilfsfunktion: Eintrag aus Datei laden
func loadEntry(entriesDir, id string) (map[string]interface{}, error) {
	entryPath := filepath.Join(entriesDir, id, id+".json")
	b, err := os.ReadFile(entryPath)
	if err != nil {
		return nil, err
	}
	var entry map[string]interface{}
	if err := json.Unmarshal(b, &entry); err != nil {
		return nil, err
	}
	return entry, nil
}

// Schnittmenge von ID-Slices
func intersectIDSets(sets [][]string) []string {
	if len(sets) == 0 {
		return nil
	}
	m := map[string]int{}
	for _, set := range sets {
		for _, id := range set {
			m[id]++
		}
	}
	var result []string
	n := len(sets)
	for id, count := range m {
		if count == n {
			result = append(result, id)
		}
	}
	return result
}
