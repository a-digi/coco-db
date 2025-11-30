package query

import (
	"encoding/json"
	"fmt"
	"github.com/a-digi/coco-db/src/table/fields"
	"github.com/a-digi/coco-db/src/index"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// FilterResult enthält die Query-Ergebnisse und die Zähler für Dateiöffnungen und RAM-Zugriffe
// Wird für API-Aggregate-Logik genutzt
//
type FilterResult struct {
	Entries      []map[string]interface{} `json:"results"`
	FileOpens    int                      `json:"fileOpens"`
	RAMHits      int                      `json:"ramHits"`
}

// Speicheroptimierte Filter-Engine: Nur Indexdaten im Speicher, sonst sequentieller Dateiscan
// Gibt die gefilterten Einträge als Array von map[string]interface{} zurück
func FilterEngine(dataDir, dbName, tableName string, query *Query, meta *fields.TableMeta) (*FilterResult, error) {
	var fileOpenCount int
	var ramHitCount int

	// Debug: Logge alle verfügbaren RAM-Index-Keys beim ersten Aufruf
	reg := index.GetRegistry()
	fmt.Printf("[DEBUG] RAM-Index-Keys: %v\n", reg.Keys())

	loadEntryCounted := func(entriesDir, id string) (map[string]interface{}, error) {
		fileOpenCount++
		return loadEntry(entriesDir, id)
	}

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

	// 1. IDs aus allen Indexfiltern sammeln
	var idSets [][]string
	for f, idxMeta := range indexedFields {
		if cond, ok := query.Filter[f]; ok {
			idxKey := dbName + "." + tableName + "." + idxMeta.Name
			reg := index.GetRegistry()
			idxObj, ok := reg.Get(idxKey)
			if ok {
				fmt.Printf("[DEBUG] RAM-Index gefunden: %s (Einträge: %d)\n", idxKey, len(idxObj))
			} else {
				fmt.Printf("[DEBUG] Kein RAM-Index für %s gefunden!\n", idxKey)
			}
			// Bereichsfilter erkennen
			switch c := cond.(type) {
			case map[string]interface{}:
				ids := filterIDsByRangeFromIndexCounted(idxObj, c, &ramHitCount)
				if len(ids) > 0 {
					idSets = append(idSets, ids)
				}
			default:
				key := fmt.Sprint(cond)
				if idsRaw, found := idxObj[key]; found {
					ramHitCount++
					if ids, ok := idsRaw.([]interface{}); ok {
						strIDs := make([]string, 0, len(ids))
						for _, id := range ids {
							if s, ok := id.(string); ok {
								strIDs = append(strIDs, s)
							}
						}
						idSets = append(idSets, strIDs)
					} else if ids, ok := idsRaw.([]string); ok {
						idSets = append(idSets, ids)
					}
				}
			}
		}
	}

	// 2. Schnittmenge aller Index-IDs bilden
	ids := intersectIDSets(idSets)

	entriesDir := filepath.Join(dataDir, dbName, tableName, "entries")
	var result []map[string]interface{}

	// 3. Wenn Index-IDs vorhanden, prüfe nur diese, sonst vollständiger Scan
	if len(ids) > 0 {
		for _, id := range ids {
			entry, err := loadEntryCounted(entriesDir, id)
			if err != nil {
				continue
			}
			// Prüfe nicht-indexierte Filter
			match := true
			for f := range nonIndexedFields {
				cond := query.Filter[f]
				val, ok := entry[f]
				if !ok {
					match = false
					break
				}
				switch c := cond.(type) {
				case map[string]interface{}:
					for op, opVal := range c {
						fn, found := operatorFuncs[op]
						if !found || !fn(val, opVal) {
							match = false
							break
						}
					}
					if !match {
						break
					}
				default:
					if !isEqual(val, c) {
						match = false
						break
					}
				}
				if !match {
					break
				}
			}
			if match {
				result = append(result, entry)
			}
		}
	} else {
		// Kein Index nutzbar: vollständiger Scan
		files, _ := os.ReadDir(entriesDir)
		for _, f := range files {
			if f.IsDir() {
				id := f.Name()
				entry, err := loadEntryCounted(entriesDir, id)
				if err != nil {
					continue
				}
				if matchesAllFiltersEngine(entry, query.Filter) {
					result = append(result, entry)
				}
			}
		}
	}

	if len(query.Sort) > 0 {
		sortEntries(result, query.Sort)
	}
	result = applyPagination(result, query.Limit, query.Offset)

	return &FilterResult{
		Entries:   result,
		FileOpens: fileOpenCount,
		RAMHits:   ramHitCount,
	}, nil
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

// Hilfsfunktion: Bereichsfilter auf Index anwenden (mit RAM-Zähler, optimiert für sortierte Keys)
func filterIDsByRangeFromIndexCounted(idxObj map[string]interface{}, cond map[string]interface{}, ramHitCount *int) []string {
	var result []string
	if len(idxObj) == 0 {
		return result
	}

	// Versuche, ob die Keys als Datum oder Zahl sortierbar sind
	var keys []string
	for k := range idxObj {
		keys = append(keys, k)
	}

	isDate := false
	isNumber := false
	if _, err := time.Parse(time.RFC3339, keys[0]); err == nil {
		isDate = true
	} else if _, err := strconv.ParseFloat(keys[0], 64); err == nil {
		isNumber = true
	}

	if isDate {
		sort.Slice(keys, func(i, j int) bool {
			t1, _ := time.Parse(time.RFC3339, keys[i])
			t2, _ := time.Parse(time.RFC3339, keys[j])
			return t1.Before(t2)
		})
	} else if isNumber {
		sort.Slice(keys, func(i, j int) bool {
			f1, _ := strconv.ParseFloat(keys[i], 64)
			f2, _ := strconv.ParseFloat(keys[j], 64)
			return f1 < f2
		})
	} else {
		sort.Strings(keys)
	}

	// Bereichsgrenzen bestimmen
	var gte, lte, gt, lt interface{}
	for op, opVal := range cond {
		switch op {
		case "gte":
			gte = opVal
		case "lte":
			lte = opVal
		case "gt":
			gt = opVal
		case "lt":
			lt = opVal
		}
	}

	for _, k := range keys {
		match := true
		if gte != nil && !compareIndexKey(k, gte, ">=", isDate) {
			match = false
		}
		if lte != nil && !compareIndexKey(k, lte, "<=", isDate) {
			match = false
		}
		if gt != nil && !compareIndexKey(k, gt, ">", isDate) {
			match = false
		}
		if lt != nil && !compareIndexKey(k, lt, "<", isDate) {
			match = false
		}
		if match {
			*ramHitCount++
			switch ids := idxObj[k].(type) {
			case []interface{}:
				for _, id := range ids {
					if s, ok := id.(string); ok {
						result = append(result, s)
					}
				}
			case []string:
				result = append(result, ids...)
			}
		}
	}
	return result
}

// Hilfsfunktion: Vergleich von Index-Keys (Datum, Zahl, String)
func compareIndexKey(key string, opVal interface{}, op string, isDate bool) bool {
	// Versuche als Zahl
	if f, err := parseFloat(key); err == nil {
		if fv, ok := toFloat(opVal); ok {
			switch op {
			case ">=":
				return f >= fv
			case "<=":
				return f <= fv
			case ">":
				return f > fv
			case "<":
				return f < fv
			case "==":
				return f == fv
			}
		}
	}
	// Versuche als Datum
	if t, err := time.Parse(time.RFC3339, key); err == nil {
		if ts, ok := opVal.(string); ok {
			if tv, err := time.Parse(time.RFC3339, ts); err == nil {
				switch op {
				case ">=":
					return !t.Before(tv)
				case "<=":
					return !t.After(tv)
				case ">":
					return t.After(tv)
				case "<":
					return t.Before(tv)
				case "==":
					return t.Equal(tv)
				}
			}
		}
	}
	// Fallback: String-Vergleich
	if sv, ok := opVal.(string); ok {
		switch op {
		case ">=":
			return key >= sv
		case "<=":
			return key <= sv
		case ">":
			return key > sv
		case "<":
			return key < sv
		case "==":
			return key == sv
		}
	}
	return false
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func toFloat(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}
