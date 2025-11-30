package query

import (
	"encoding/json"
	"fmt"
	"github.com/a-digi/coco-db/src/query/filter"
	"github.com/a-digi/coco-db/src/table/fields"
	"os"
	"path/filepath"
	"sort"
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
			idxPath := filepath.Join(dataDir, dbName, tableName, "indexes", fields.GetIndexFileName(idxMeta.Name))
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
	if query.IsSearchQuery {
		// Eigene Filter-Logik für SearchQuery
		files, _ := os.ReadDir(entriesDir)
		for _, f := range files {
			if f.IsDir() {
				id := f.Name()
				entry, err := loadEntry(entriesDir, id)
				if err == nil && filter.MatchesAllFiltersSearch(entry, query.Filter, operatorFuncs, isEqual) {
					fmt.Printf("[JOIN-LOG] Join-Treffer: id=%v, Filter=%v\n", id, query.Filter)
					result = append(result, entry)
				}
			}
		}
	} else {
		if len(ids) > 0 {
			for _, id := range ids {
				entry, err := loadEntry(entriesDir, id)
				if err == nil && matchesAllFiltersEngine(entry, query.Filter) {
					result = append(result, entry)
				}
			}
		} else {
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
	}

	// Nach dem Filtern: Sortierung und Paginierung
	if len(query.Sort) > 0 {
		sortEntries(result, query.Sort)
	}
	result = applyPagination(result, query.Limit, query.Offset)
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

// Hilfsfunktion: IDs aus Indexdatei laden (echtes Indexformat: map[string][]string)
func loadIDsFromIndex(idxPath string, cond interface{}) ([]string, error) {
	b, err := os.ReadFile(idxPath)
	if err != nil {
		return nil, err
	}
	// Versuche zuerst map[string][]string
	var idxObj map[string][]string
	if err := json.Unmarshal(b, &idxObj); err == nil {
		key := fmt.Sprint(cond)
		if ids, ok := idxObj[key]; ok {
			return ids, nil
		}
		return []string{}, nil
	}
	// Fallback: []string (ältere Tests)
	var idxArr []string
	if err := json.Unmarshal(b, &idxArr); err == nil {
		return idxArr, nil
	}
	return nil, fmt.Errorf("Indexdatei hat unbekanntes Format")
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

// sortEntries sortiert die Einträge nach den angegebenen Feldern (auf- und absteigend)
func sortEntries(entries []map[string]interface{}, sortFields []string) {
	sort.SliceStable(entries, func(i, j int) bool {
		for _, field := range sortFields {
			asc := true
			name := field
			if len(field) > 0 && field[0] == '-' {
				asc = false
				name = field[1:]
			}
			vi, iok := entries[i][name]
			vj, jok := entries[j][name]
			if !iok || !jok {
				continue
			}
			cmp := compareValues(vi, vj)
			if cmp == 0 {
				continue
			}
			if asc {
				return cmp < 0
			} else {
				return cmp > 0
			}
		}
		return false
	})
}

// compareValues vergleicht zwei Werte (int, float, string, bool)
func compareValues(a, b interface{}) int {
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	if okA && okB {
		if fa < fb {
			return -1
		} else if fa > fb {
			return 1
		}
		return 0
	}
	sa, okA := a.(string)
	sb, okB := b.(string)
	if okA && okB {
		if sa < sb {
			return -1
		} else if sa > sb {
			return 1
		}
		return 0
	}
	ba, okA := a.(bool)
	bb, okB := b.(bool)
	if okA && okB {
		if ba == bb {
			return 0
		} else if !ba && bb {
			return -1
		} else {
			return 1
		}
	}
	return 0
}

// applyPagination schneidet das Ergebnis auf limit/offset zu
func applyPagination(entries []map[string]interface{}, limit, offset int) []map[string]interface{} {
	if offset > len(entries) {
		return []map[string]interface{}{}
	}
	end := len(entries)
	if limit > 0 && offset+limit < end {
		end = offset + limit
	}
	return entries[offset:end]
}
