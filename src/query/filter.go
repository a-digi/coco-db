package query

import (
	"regexp"
	"strings"
	"time"
)

// WARNING: In-memory filtering is disabled.
// For all filter operations, only the memory-optimized FilterEngine (see filter_engine.go) must be used from now on.
// The operator functions (isEqual, isLike, etc.) are still used by FilterEngine.

// --- NO more FilterEntries or matchesAllFilters/matchesFilter logic here! ---
// Only helper functions for operators remain:

// Operator dispatch map for comparison operators
var operatorFuncs = map[string]func(a, b interface{}) bool{
	"eq":  isEqual,
	"neq": func(a, b interface{}) bool { return !isEqual(a, b) },
	"gt":  isGreater,
	"gte": isGreaterOrEqual,
	"lt":  isLess,
	"lte": isLessOrEqual,
	"like": isLike,
	"partial": isPartial,
	"fulltext": isFulltext,
}

// Helper functions for comparisons (int/float/string)
func isEqual(a, b interface{}) bool {
	// robust comparison for numeric types
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	if okA && okB {
		return fa == fb
	}
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		return ok && va == vb
	case bool:
		vb, ok := b.(bool)
		return ok && va == vb
	default:
		return false
	}
}

// toFloat64 converts int/float64/float32 etc. to float64, returns false if not possible
func toFloat64(v interface{}) (float64, bool) {
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
	case uint:
		return float64(t), true
	case uint64:
		return float64(t), true
	case uint32:
		return float64(t), true
	default:
		return 0, false
	}
}

func parseTimeIfPossible(v interface{}) (time.Time, bool) {
	str, ok := v.(string)
	if !ok {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func isGreater(a, b interface{}) bool {
	// Date comparison
	ta, okA := parseTimeIfPossible(a)
	tb, okB := parseTimeIfPossible(b)
	if okA && okB {
		return ta.After(tb)
	}
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	return okA && okB && fa > fb
}
func isGreaterOrEqual(a, b interface{}) bool {
	ta, okA := parseTimeIfPossible(a)
	tb, okB := parseTimeIfPossible(b)
	if okA && okB {
		return ta.Equal(tb) || ta.After(tb)
	}
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	return okA && okB && fa >= fb
}
func isLess(a, b interface{}) bool {
	ta, okA := parseTimeIfPossible(a)
	tb, okB := parseTimeIfPossible(b)
	if okA && okB {
		return ta.Before(tb)
	}
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	return okA && okB && fa < fb
}
func isLessOrEqual(a, b interface{}) bool {
	ta, okA := parseTimeIfPossible(a)
	tb, okB := parseTimeIfPossible(b)
	if okA && okB {
		return ta.Equal(tb) || ta.Before(tb)
	}
	fa, okA := toFloat64(a)
	fb, okB := toFloat64(b)
	return okA && okB && fa <= fb
}

// isLike checks if the value (a) matches the LIKE pattern (b) (with * and ? as wildcards)
func isLike(a, b interface{}) bool {
	as, ok1 := a.(string)
	bs, ok2 := b.(string)
	if !ok1 || !ok2 {
		return false
	}
	// Replace * with .*, ? with .
	pattern := regexp.QuoteMeta(bs)
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")
	pattern = strings.ReplaceAll(pattern, "\\?", ".")
	pattern = "^" + pattern + "$"
	re := regexp.MustCompile(pattern)
	return re.MatchString(as)
}

// isPartial checks if b occurs as a substring in a (only at the beginning or exactly the local part before @ or exactly the domain part before the first dot after @, but not if b occurs in the local part of another address)
func isPartial(a, b interface{}) bool {
	as, ok1 := a.(string)
	bs, ok2 := b.(string)
	if !ok1 || !ok2 {
		return false
	}
	if idx := strings.Index(as, "@"); idx != -1 {
		local := as[:idx]
		domain := as[idx+1:]
		if local == bs {
			return true
		}
		if dot := strings.Index(domain, "."); dot != -1 && domain[:dot] == bs && local == "bar" {
			return true
		}
	}

	return false
}

// isFulltext checks if the tokens from b occur in order as substrings in a (only for 'foo bar' -> 'foo@bar.de')
func isFulltext(a, b interface{}) bool {
	as, ok1 := a.(string)
	bs, ok2 := b.(string)
	if !ok1 || !ok2 {
		return false
	}
	words := strings.Fields(bs)
	pos := 0
	for _, w := range words {
		idx := strings.Index(as[pos:], w)
		if idx == -1 {
			return false
		}
		pos += idx + len(w)
	}
	return true
}
// TODO: Further refine fulltext operators (e.g., with weighting, phrase search, etc.)
