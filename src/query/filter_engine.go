package query

import (
	"encoding/json"
	"fmt"
	"github.com/a-digi/coco-db/src/table/fields"
	"github.com/a-digi/coco-db/src/index"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

// FilterResult contains the query results and the counters for file openings and RAM accesses
// Used for API aggregate logic
type FilterResult struct {
	Entries      []map[string]interface{} `json:"results"`
	FileOpens    int                      `json:"fileOpens"`
	RAMHits      int                      `json:"ramHits"`
}

// Memory-optimized filter engine: Only index data in memory, otherwise sequential file scan
// Returns the filtered entries as an array of map[string]interface{}
func FilterEngine(dataDir, dbName, tableName string, query *Query, meta *fields.TableMeta) (*FilterResult, error) {
	var fileOpenCount int
	var ramHitCount int
    // Directory where entries are stored
	entriesDir := filepath.Join(dataDir, dbName, tableName, "entries")
	var result []map[string]interface{}

	indexedFields := buildIndexedFields(meta)

    // Store ID sets from each index filter
	var idSets [][]string
    // 1. Find IDs from RAM index for each indexed field in the filter
	idSets = findInIndexRam(dbName, tableName, indexedFields, query, &ramHitCount)

	// 2. Form the intersection of all index IDs
	ids := intersectIDSets(idSets)

	if len(idSets) > 0 && len(ids) == 0 {
		fmt.Printf("[FILTER-TRACE] Index IDs found, but intersection is empty. No entries will be loaded.\n")
		return &FilterResult{Entries: []map[string]interface{}{}, FileOpens: fileOpenCount, RAMHits: ramHitCount}, nil
	}

	// 3. If index IDs are present, only check these, otherwise full scan
	if len(ids) > 0 {
		fmt.Printf("[FILTER-TRACE] Using index IDs (%d) for loading entries.\n", len(ids))
		for _, id := range ids {
			entry, err := loadEntryCounted(entriesDir, id, &fileOpenCount)
			if err != nil {
				continue
			}
            // Final check against all filter conditions
			match := matcheFilters(entry, query.Filter)

			if match {
				result = append(result, entry)
			}
		}
	} else if inIDs := getInFilterIDs(query.Filter); len(inIDs) > 0 {
		fmt.Printf("[FILTER-TRACE] Using IN-filter IDs (%d) for loading entries.\n", len(inIDs))
		for _, id := range inIDs {
			entry, err := loadEntryCounted(entriesDir, id, &fileOpenCount)
			if err != nil {
				continue
			}
			if matcheFilters(entry, query.Filter) {
				result = append(result, entry)
			}
		}
	} else {
		fmt.Printf("[FILTER-TRACE] No index or IN-filter found. Performing full scan of all entries!\n")
		files, _ := os.ReadDir(entriesDir)
		for _, f := range files {
			if f.IsDir() {
				id := f.Name()
				entry, err := loadEntryCounted(entriesDir, id, &fileOpenCount)
				if err != nil {
					continue
				}
				if matcheFilters(entry, query.Filter) {
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

// buildIndexedFields extracts all indexed fields from the TableMeta
func buildIndexedFields(meta *fields.TableMeta) map[string]fields.IndexMeta {
	indexedFields := map[string]fields.IndexMeta{}

	for _, idx := range meta.Indexes {
		for _, f := range idx.Fields {
			indexedFields[f] = idx
		}
	}

	return indexedFields
}

// matcheFilters checks if an entry matches all filter conditions
func matcheFilters(entry map[string]interface{}, filter map[string]interface{}) bool {
	for f := range filter {
		cond := filter[f]
		val, ok := entry[f]
		if !ok {
			return false
		}
		switch c := cond.(type) {
		case map[string]interface{}:
			for op, opVal := range c {
				if op == "like" {
					valStr, ok1 := val.(string)
					pattern, ok2 := opVal.(string)
					if !ok1 || !ok2 || !matchLikePattern(valStr, pattern) {
						return false
					}
					continue
				}
				fn, found := operatorFuncs[op]
				if !found || !fn(val, opVal) {
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

// matchLikePattern checks if val meets the LIKE pattern (* and ? as wildcards)
func matchLikePattern(val, pattern string) bool {
	// Replace * with .* and ? with . for regular expression
	regex := "^" + pattern + "$"
	regex = regexp.MustCompile(`([\\.\\+\\[\\]\\(\\)\\^\\$\\|\\{\\}])`).ReplaceAllString(regex, `\\$1`)
	regex = regexp.MustCompile(`\\*`).ReplaceAllString(regex, ".*")
	regex = regexp.MustCompile(`\\?`).ReplaceAllString(regex, ".")
	matched, err := regexp.MatchString(regex, val)
	return err == nil && matched
}

// Helper function: Load IDs from index file (actual index format: map[string][]string)
func loadIDsFromIndex(idxPath string, cond interface{}) ([]string, error) {
	b, err := os.ReadFile(idxPath)
	if err != nil {
		return nil, err
	}
	// First try map[string][]string
	var idxObj map[string][]string
	if err := json.Unmarshal(b, &idxObj); err == nil {
		key := fmt.Sprint(cond)
		if ids, ok := idxObj[key]; ok {
			return ids, nil
		}
		return []string{}, nil
	}
	// Fallback: []string (older tests)
	var idxArr []string
	if err := json.Unmarshal(b, &idxArr); err == nil {
		return idxArr, nil
	}
	return nil, fmt.Errorf("Index file has unknown format")
}

// Helper function: Load entry from file
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

// loadEntryCounted loads an entry and increments the file open counter
func loadEntryCounted(entriesDir, id string, fileOpenCount *int) (map[string]interface{}, error) {
	*fileOpenCount++
	return loadEntry(entriesDir, id)
}

// Intersection of ID slices
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

// sortEntries sorts the entries by the specified fields (ascending and descending)
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

// compareValues compares two values (int, float, string, bool)
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

// applyPagination trims the result to limit/offset
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

// Helper function: Apply range filter to index (with RAM counter, optimized for sorted keys)
func filterIDsByRangeFromIndexCounted(idxObj map[string]interface{}, cond map[string]interface{}, ramHitCount *int) []string {
	var result []string

	if len(idxObj) == 0 {
		return result
	}

	// Try if the keys can be sorted as date or number
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

	// Determine range boundaries
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

// Helper function: Compare index keys (date, number, string)
func compareIndexKey(key string, opVal interface{}, op string, isDate bool) bool {
	// Try as number
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
	// Try as date
	if t, err := time.Parse(time.RFC3339, key); err == nil {
		if ts, ok := opVal.(string); ok {
			if tv, err := time.Parse(time.RFC3339, ts); err == nil {
				// Explicitly normalize both values to UTC
				t = t.UTC()
				tv = tv.UTC()
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

	// Fallback: String comparison
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

// Helper function: Extract IDs from an "in" filter
func getInFilterIDs(filter map[string]interface{}) []string {
    for _, cond := range filter {
        if condMap, ok := cond.(map[string]interface{}); ok {
            if inVal, ok := condMap["in"]; ok {
                var ids []string
                switch v := inVal.(type) {
                case []interface{}:
                    for _, id := range v {
                        if s, ok := id.(string); ok {
                            ids = append(ids, s)
                        }
                    }
                case []string:
                    ids = append(ids, v...)
                }
                return ids
            }
        }
    }
    return nil
}

// findInIndexRam searches IDs in the RAM index for a field and returns the ID slices
func findInIndexRam(dbName, tableName string, indexedFields map[string]fields.IndexMeta, query *Query, ramHitCount *int) [][]string {
	var idSets [][]string

	// For each indexed field in the filter, get IDs from RAM index
	for f, idxMeta := range indexedFields {
		cond, ok := query.Filter[f]

		if !ok {
			continue
		}

		idxKey := dbName + "." + tableName + "." + idxMeta.Name
		reg := index.GetRegistry()
		idxObj, ok := reg.Get(idxKey)

		if !ok {
			continue
		}

		// Detect range filters
		switch c := cond.(type) {
		case map[string]interface{}:
			ids := filterIDsByRangeFromIndexCounted(idxObj, c, ramHitCount)
			if len(ids) > 0 {
				idSets = append(idSets, ids)
			}
		default:
			key := fmt.Sprint(cond)
			idsRaw, found := idxObj[key]
			if !found {
				continue
			}
			(*ramHitCount)++

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

	return idSets
}
